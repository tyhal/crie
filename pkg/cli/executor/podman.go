package executor

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime/trace"
	"strconv"
	"sync"
	"time"

	"charm.land/log/v2"
	"github.com/containerd/platforms"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.podman.io/podman/v6/pkg/api/handlers"
	"go.podman.io/podman/v6/pkg/bindings"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/bindings/system"
	"go.podman.io/podman/v6/pkg/machine/define"
	"go.podman.io/podman/v6/pkg/machine/env"
	"go.podman.io/podman/v6/pkg/machine/provider"
	"go.podman.io/podman/v6/pkg/machine/vmconfigs"
	"go.podman.io/podman/v6/pkg/specgen"

	"github.com/tyhal/crie/pkg/linter"
)

// podmanExecutor runs CLI tools inside a Podman container.
type podmanExecutor struct {
	Instance
	image string
	// conn is only used for its values, the bindings keep the connection under unexported context keys
	conn context.Context
	id   string
}

// podmanCtx layers the connection's values beneath a call's ctx, so the call keeps ctx's cancellation and span
type podmanCtx struct {
	context.Context
	conn context.Context
}

func (c podmanCtx) Value(key any) any {
	if v := c.Context.Value(key); v != nil {
		return v
	}
	return c.conn.Value(key)
}

// with returns ctx able to make podman binding calls on this executor's connection
func (e *podmanExecutor) with(ctx context.Context) context.Context {
	return podmanCtx{ctx, e.conn}
}

// NewPodman creates an executor using containers managed by the podman client.
func NewPodman(image string) Executor {
	return &podmanExecutor{
		image: image,
	}
}

var podmanImagePullLocks sync.Map

// podmanHost is what probing podman learns, it can't change during a run
type podmanHost struct {
	uri      string
	rootless bool
}

// cachedPodman probes podman once for every linter's WillPodman and Setup,
// without a caller's context so one aborted caller can't cache a failure for everyone
var cachedPodman = sync.OnceValues(func() (podmanHost, error) {
	return probePodman(context.Background())
})

// WillPodman checks whether Podman is available and responsive on the host.
// If Podman is configured (a machine or CONTAINER_HOST exists) but doesn't respond, the error wraps ErrRuntimeUnavailable.
func WillPodman(_ context.Context) error {
	_, err := cachedPodman()
	return err
}

func probePodman(ctx context.Context) (podmanHost, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sock, err := getPodmanMachineSocket()
	if err != nil {
		return podmanHost{}, err
	}

	unavailable := func(err error) error {
		if sock.configured {
			return fmt.Errorf("%w: podman is configured but not responding (is the machine running?): %w", ErrRuntimeUnavailable, err)
		}
		return err
	}

	c, err := bindings.NewConnection(ctx, sock.uri)
	if err != nil {
		return podmanHost{}, unavailable(err)
	}

	info, err := system.Info(c, &system.InfoOptions{})
	if err != nil {
		return podmanHost{}, unavailable(err)
	}

	return podmanHost{uri: sock.uri, rootless: info.Host.Security.Rootless}, nil
}

// podmanSocket is where to reach podman
type podmanSocket struct {
	uri string
	// configured is set when the user has set podman up, rather than uri being a guessed default
	configured bool
}

// getPodmanMachineSocket finds the podman socket, looking for the default machine under every provider
// as it may have been created with one other than the current default (e.g. applehv rather than libkrun)
func getPodmanMachineSocket() (podmanSocket, error) {
	if _, found := os.LookupEnv("CONTAINER_HOST"); found {
		return podmanSocket{configured: true}, nil
	}

	providers := provider.GetAll()
	if p, err := provider.Get(); err == nil {
		providers = append([]vmconfigs.VMProvider{p}, providers...)
	}
	for _, p := range providers {
		dirs, err := env.GetMachineDirs(p.VMType())
		if err != nil {
			continue
		}
		mc, err := vmconfigs.LoadMachineByName(define.DefaultMachineName, dirs)
		if err != nil {
			continue
		}
		sock, _, err := mc.ConnectionInfo(p.VMType())
		if err != nil {
			return podmanSocket{}, err
		}
		return podmanSocket{uri: "unix://" + sock.Path, configured: true}, nil
	}

	// no machine, so expect a native rootless podman. It's not counted as configured even if its socket file exists,
	// as stopping podman.socket leaves the file behind (systemd doesn't remove it), which shouldn't stop a fallback to docker
	currentUser, err := user.Current()
	if err != nil {
		return podmanSocket{}, err
	}
	return podmanSocket{uri: fmt.Sprintf("unix:///run/user/%s/podman/podman.sock", currentUser.Uid)}, nil
}

// Setup creates and starts a disposable Podman container for executing commands.
func (e *podmanExecutor) Setup(ctx context.Context, i Instance) error {
	e.Instance = i

	host, err := cachedPodman()
	if err != nil {
		return fmt.Errorf("probing podman: %w", err)
	}

	c, err := bindings.NewConnection(ctx, host.uri)
	if err != nil {
		return fmt.Errorf("creating podman client: %w", err)
	}
	// trace API calls the way the docker client does
	if conn, err := bindings.GetClient(c); err == nil {
		conn.Client.Transport = otelhttp.NewTransport(conn.Client.Transport, otelhttp.WithSpanNameFormatter(apiSpanName))
	}
	e.conn = c
	ctx = e.with(ctx)

	if err := e.pull(ctx); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	b := make([]byte, 4)
	rand.Read(b)
	shortid := hex.EncodeToString(b)

	wdHost, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	wdContainer, err := GetWorkdirAsLinuxPath()
	if err != nil {
		return fmt.Errorf("getting containerWorkdir as a linux path: %w", err)
	}
	cacheContainer := "/tmp/crie_cache"

	currPlatform := platforms.DefaultSpec()
	currPlatform.OS = "linux"

	s := specgen.NewSpecGenerator(e.image, false)
	s.Name = fmt.Sprintf("crie-%s-%s", filepath.Base(e.Bin), shortid)
	s.Entrypoint = []string{"/bin/sh", "-c"}
	s.Command = []string{"tail -f /dev/null"}
	s.Env = map[string]string{
		"XDG_CACHE_HOME": cacheContainer,
	}
	s.WorkDir = wdContainer
	s.NetNS = specgen.Namespace{
		NSMode: specgen.NoNetwork,
	}

	if host.rootless {
		// root in container should map to the host user in rootless mode
		s.User = "0"
	} else {
		// avoid using root on a rootfull podman
		s.User = strconv.Itoa(os.Getuid())
	}

	mountPerms := "ro" // should be used for static-only
	if e.WillWrite {
		mountPerms = "rw" // should be used for formatters
	}

	// U changes ownership, breaks on Linux but might have been used previously as a workaround on Mac
	s.Mounts = []spec.Mount{
		{
			Type:        "bind",
			Source:      wdHost,
			Destination: wdContainer,
			Options: []string{
				"rbind",
				mountPerms,
				// "z" instead of "Z" because we run multiple containers in parallel
				"z",
			},
		},
	}
	s.Volumes = []*specgen.NamedVolume{
		{
			Name: "crie-cache",
			Dest: cacheContainer,
			Options: []string{
				"z",
			},
		},
	}

	log.With(
		"image", e.image,
		"workdirHost", wdHost,
		"workdirContainer", wdContainer,
		"cacheContainer", cacheContainer,
		"canWrite", e.WillWrite).Debug("starting container")

	createResponse, err := containers.CreateWithSpec(ctx, s, nil)
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}
	e.id = createResponse.ID

	startOptions := containers.StartOptions{}
	if err := containers.Start(ctx, e.id, &startOptions); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	return nil
}

func (e *podmanExecutor) pull(ctx context.Context) error {
	lock, _ := podmanImagePullLocks.LoadOrStore(e.image, &sync.Mutex{})
	mu := lock.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	exists, err := images.Exists(ctx, e.image, &images.ExistsOptions{})
	if err != nil {
		return fmt.Errorf("failed to check if image exists: %v", err)
	}
	if !exists {
		if _, err = images.Pull(ctx, e.image, nil); err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}
	}

	return nil
}

// Exec runs the configured command inside the prepared Podman container.
func (e *podmanExecutor) Exec(ctx context.Context, stdout io.Writer, stderr io.Writer, targets ...string) error {
	ctx = e.with(ctx)
	defer trace.StartRegion(ctx, "Exec "+e.Bin).End()

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

	// disabled on rootless/linux
	//currentUser, err := user.Current()
	//if err != nil {
	//	return fmt.Errorf("getting current user: %w", err)
	//}

	execCreateConfig := handlers.ExecCreateConfig{
		ExecCreateRequest: container.ExecCreateRequest{
			Privileged:   false,
			Tty:          false,
			ConsoleSize:  nil,
			AttachStdin:  false,
			AttachStderr: true,
			AttachStdout: true,
			WorkingDir:   wdContainer,
			Env:          traceEnv(ctx),
			Cmd:          cmd,
		},
	}

	execID, err := containers.ExecCreate(ctx, e.id, &execCreateConfig)
	if err != nil {
		return fmt.Errorf("creating exec: %w", err)
	}
	defer func() {
		// removed even if this run was cancelled
		err = containers.ExecRemove(context.WithoutCancel(ctx), execID, &containers.ExecRemoveOptions{})
		if err != nil && ctx.Err() == nil {
			log.Error(fmt.Errorf("closing exec: %w", err))
		}
	}()

	logs, err := attachedExecStart(ctx, execID, &containers.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("attaching to created exec: %w", err)
	}
	defer func() {
		err := logs.Close()
		if err != nil {
			log.Error(fmt.Errorf("closing logs: %w", err))
		}
	}()

	_, err = stdcopy.StdCopy(stdout, stderr, logs)
	if err != nil {
		// an aborted run cuts the stream, which isn't worth reporting
		if ctx.Err() != nil {
			return context.Cause(ctx)
		}
		log.Errorf("Error demultiplexing logs: %v", err)
	}

	inspect, err := containers.ExecInspect(ctx, execID, &containers.ExecInspectOptions{})
	if err != nil {
		return fmt.Errorf("inspecting exec: %w", err)
	}
	if inspect.Running {
		return errors.New("container still running after end of attach output stream")
	}
	if inspect.ExitCode != 0 {
		return linter.Result(errors.New("exit code " + strconv.Itoa(inspect.ExitCode)))
	}
	return nil
}

// Cleanup kills and removes the temporary Podman container created during Setup.
// It's killed rather than stopped as it's only idling once the tools exec'd into it have finished,
// and a graceful stop waits out its timeout for any image whose PID 1 ignores SIGTERM, like our /bin/sh.
func (e *podmanExecutor) Cleanup(ctx context.Context) error {

	// TODO cleanup based on labels (project, language)

	if e.id != "" {
		log.With("podmanId", e.id).Debug("removing container")
		removeOptions := containers.RemoveOptions{
			Force:   new(true),
			Timeout: new(uint(0)),
		}
		if _, err := containers.Remove(e.with(ctx), e.id, &removeOptions); err != nil {
			return err
		}
		e.id = ""
	}

	return nil
}

// modified to capture output from containers.ExecStart
func attachedExecStart(ctx context.Context, sessionID string, options *containers.ExecStartOptions) (io.ReadCloser, error) {
	if options == nil {
		options = new(containers.ExecStartOptions)
	}
	_ = options
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	body := struct {
		Detach bool `json:"Detach" yaml:"Detach"`
	}{
		Detach: false,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := conn.DoRequest(ctx, bytes.NewReader(bodyJSON), http.MethodPost, "/exec/%s/start", nil, nil, sessionID)
	if resp == nil || err != nil {
		return nil, err
	}

	return resp.Body, err
}
