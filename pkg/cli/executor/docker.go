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

	"github.com/tyhal/crie/pkg/linter"
)

// dockerExecutor runs CLI tools inside a Docker container.
type dockerExecutor struct {
	Instance
	image      string
	client     *client.Client
	execCtx    context.Context
	execCancel context.CancelFunc
	id         string
}

// NewDocker creates an executor using containers managed by the docker client.
func NewDocker(image string) Executor {
	return &dockerExecutor{
		image: image,
	}
}

var dockerInstalled = false
var dockerImagePullLocks sync.Map

// WillDocker checks whether Docker is available on the host (via the socket).
func WillDocker() error {
	if dockerInstalled {
		return nil
	}
	_, err := os.Stat("/var/run/docker.sock")
	if err != nil {
		return err
	}
	dockerInstalled = true
	return nil
}

// Setup creates and starts a disposable Docker container for executing commands.
func (e *dockerExecutor) Setup(ctx context.Context, i Instance) error {
	e.Instance = i

	// Add our clientDocker
	{
		c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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

	e.execCtx, e.execCancel = context.WithCancel(ctx)
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
	_, err = io.Copy(&pullOut, pullStat)
	if log.DebugLevel >= log.GetLevel() {
		fmt.Print(pullOut.String())
	}

	return nil
}

// Exec runs the configured command inside the prepared Docker container.
func (e *dockerExecutor) Exec(filePath string, stdout io.Writer, stderr io.Writer) error {

	// working solution posted to https://stackoverflow.com/questions/52145231/cannot-get-logs-from-docker-container-using-golang-docker-sdk

	// Ensure we can mount our filesystem to the same path inside the container
	targetFile := ToLinuxPath(filePath)
	wdContainer, err := GetWorkdirAsLinuxPath()
	if err != nil {
		return err
	}

	if e.ChDir {
		wdContainer = filepath.Join(wdContainer, filepath.Dir(targetFile))
		targetFile = filepath.Base(targetFile)
	}

	cmd := make([]string, 0, 1+len(e.Start)+1+len(e.End))
	cmd = append([]string{e.Bin}, e.Start...)
	cmd = append(cmd, targetFile)
	cmd = append(cmd, e.End...)

	log.Debug(cmd)
	config := container.ExecOptions{
		Cmd:          cmd,
		User:         fmt.Sprintf("%d", os.Getuid()),
		WorkingDir:   wdContainer,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          false,
	}
	execResp, err := e.client.ContainerExecCreate(e.execCtx, e.id, config)
	if err != nil {
		return err
	}

	startConfig := container.ExecAttachOptions{
		Detach: false,
		Tty:    false,
	}
	attach, err := e.client.ContainerExecAttach(e.execCtx, execResp.ID, startConfig)
	if err != nil {
		return err
	}
	defer attach.Close()

	_, err = stdcopy.StdCopy(stdout, stderr, attach.Reader)
	if err != nil {
		log.Errorf("Error demultiplexing logs: %v", err)
	}

	inspect, err := e.client.ContainerExecInspect(e.execCtx, execResp.ID)
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

// Cleanup stops and removes the temporary Docker container created during Setup.
func (e *dockerExecutor) Cleanup(ctx context.Context) error {

	if e.execCancel != nil {
		defer e.execCancel()
	}

	if e.id != "" {
		d := log.With("dockerId", e.id)

		d.Debug("stopping container")
		if err := e.client.ContainerStop(ctx, e.id, container.StopOptions{Timeout: new(1)}); err != nil {
			return err
		}
		d.Debug("removing container")
		if err := e.client.ContainerRemove(ctx, e.id, container.RemoveOptions{}); err != nil {
			return err
		}
		e.id = ""
	}

	return nil
}
