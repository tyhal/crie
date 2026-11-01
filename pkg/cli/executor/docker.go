package executor

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	log "charm.land/log/v2"
	"github.com/containerd/platforms"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/tyhal/crie/pkg/linter"
)

// dockerExecutor runs CLI tools inside a Docker container.
type dockerExecutor struct {
	Instance
	image  string
	client *client.Client
	id     string
}

// NewDocker creates an executor using containers managed by the docker client.
func NewDocker(image string) Executor {
	return &dockerExecutor{
		image: image,
	}
}

var dockerImagePullLocks sync.Map

// dockerSocket caches the socket check, it's called by every linter's setup
var dockerSocket = sync.OnceValue(func() error {
	_, err := os.Stat("/var/run/docker.sock")
	return err
})

// WillDocker checks whether Docker is available on the host (via the socket).
func WillDocker() error {
	return dockerSocket()
}

// Setup creates and starts a disposable Docker container for executing commands.
func (e *dockerExecutor) Setup(ctx context.Context, i Instance) error {
	e.Instance = i

	// Add our clientDocker
	{
		c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(),
			client.WithTraceOptions(otelhttp.WithSpanNameFormatter(apiSpanName)))
		if err != nil {
			return fmt.Errorf("failed to create docker client: %w", err)
		}
		e.client = c
	}

	if err := e.pull(ctx); err != nil {
		return err
	}

	// TODO do these both in some common container helpers
	wdHost, err := os.Getwd()
	if err != nil {
		return err
	}
	wdContainer, err := GetWorkdirAsLinuxPath()
	if err != nil {
		return err
	}
	cacheContainer := "/tmp/crie_cache"

	currPlatform := platforms.DefaultSpec()
	currPlatform.OS = "linux"

	b := make([]byte, 4)
	_, _ = rand.Read(b)
	shortid := hex.EncodeToString(b)

	labels := map[string]string{
		"owner": "crie",
	}

	resp, err := e.client.ContainerCreate(ctx,
		&container.Config{
			Entrypoint: []string{},
			Cmd:        []string{"/bin/sh", "-c", "tail -f /dev/null"},
			Env: []string{
				"XDG_CACHE_HOME=" + cacheContainer,
			},
			Labels:          labels,
			Image:           e.image,
			WorkingDir:      wdContainer,
			NetworkDisabled: true,
			User:            fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:     mount.TypeBind,
					Source:   wdHost,
					Target:   wdContainer,
					ReadOnly: !e.WillWrite,
				},
				{
					Type:   mount.TypeVolume,
					Source: "crie-cache",
					Target: cacheContainer,
					VolumeOptions: &mount.VolumeOptions{
						Labels: labels,
						DriverConfig: &mount.Driver{
							Name: "local",
							Options: map[string]string{
								"type":   "tmpfs",
								"device": "tmpfs",
								"o":      fmt.Sprintf("uid=%d", os.Getuid()),
							},
						},
					},
				},
			},
		},
		nil,
		&currPlatform,
		fmt.Sprintf("crie-%s-%s", filepath.Base(e.Bin), shortid),
	)
	if err != nil {
		return err
	}
	e.id = resp.ID

	return e.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

func (e *dockerExecutor) pull(ctx context.Context) error {
	lock, _ := dockerImagePullLocks.LoadOrStore(e.image, &sync.Mutex{})
	mu := lock.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	_, err := e.client.ImageHistory(ctx, e.image)
	// exit early if the image is already present
	if err == nil {
		return nil
	}

	pullStat, err := e.client.ImagePull(ctx, e.image, image.PullOptions{})
	if err != nil {
		log.With("stage", "docker pull", "image", e.image).Fatal(err)
		return err
	}

	var pullOut bytes.Buffer
	_, _ = io.Copy(&pullOut, pullStat)
	if log.DebugLevel >= log.GetLevel() {
		fmt.Print(pullOut.String())
	}

	return nil
}

// Exec runs the configured command inside the prepared Docker container.
func (e *dockerExecutor) Exec(ctx context.Context, stdout io.Writer, stderr io.Writer, targets ...string) error {
	if len(targets) == 0 {
		return ErrNoTargets
	}
	targetFiles := make([]string, len(targets))
	for i, target := range targets {
		targetFiles[i] = ToLinuxPath(target)
	}
	wdContainer, err := containerWorkdir(targets[0], e.ChDir, e.NoFileArg)
	if err != nil {
		return err
	}

	cmd := append([]string{e.Bin}, e.buildParams(targetFiles...)...)

	log.Debug(cmd)
	config := container.ExecOptions{
		Cmd:          cmd,
		User:         fmt.Sprintf("%d", os.Getuid()),
		WorkingDir:   wdContainer,
		Env:          traceEnv(ctx),
		AttachStderr: true,
		AttachStdout: true,
		Tty:          false,
	}
	execResp, err := e.client.ContainerExecCreate(ctx, e.id, config)
	if err != nil {
		return err
	}

	startConfig := container.ExecAttachOptions{
		Detach: false,
		Tty:    false,
	}
	attach, err := e.client.ContainerExecAttach(ctx, execResp.ID, startConfig)
	if err != nil {
		return err
	}
	defer attach.Close()

	_, err = stdcopy.StdCopy(stdout, stderr, attach.Reader)
	if err != nil {
		// an aborted run cuts the stream, which isn't worth reporting
		if ctx.Err() != nil {
			return context.Cause(ctx)
		}
		log.Errorf("Error demultiplexing logs: %v", err)
	}

	inspect, err := e.client.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return err
	}
	if inspect.Running {
		return errors.New("container still running after end of attach output stream")
	}
	if inspect.ExitCode != 0 {
		return linter.Result(fmt.Errorf("exit code %d", inspect.ExitCode))
	}
	return nil
}

// Cleanup kills and removes the temporary Docker container created during Setup.
// It's killed rather than stopped as it's only idling once the tools exec'd into it have finished,
// and a graceful stop waits out its timeout for any image whose PID 1 ignores SIGTERM, like our /bin/sh.
func (e *dockerExecutor) Cleanup(ctx context.Context) error {

	if e.id != "" {
		log.With("dockerId", e.id).Debug("removing container")
		if err := e.client.ContainerRemove(ctx, e.id, container.RemoveOptions{Force: true}); err != nil {
			return err
		}
		e.id = ""
	}

	return nil
}
