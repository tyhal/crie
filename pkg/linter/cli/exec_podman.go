package cli

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/containerd/platforms"
	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/docker/docker/api/types/container"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

var podmanInstalled = false

func (e *Lint) willPodman() error {

	if e.Container.Image == "" {
		return errors.New("no image specified for configuration " + e.Name() + "")
	}
	if podmanInstalled {
		return nil
	}

	_, err := exec.LookPath("podman")
	if err != nil {
		return err
	}

	podmanInstalled = true
	return nil
}

func (e *Lint) startPodman() error {
	ctx := context.Background()

	{
		currentUser, err := user.Current()
		if err != nil {
			return err
		}

		c, err := bindings.NewConnection(ctx, fmt.Sprintf("unix:///run/user/%s/podman/podman.sock", currentUser.Uid))
		if err != nil {
			return err
		}
		e.Container.clientPodman = &c
	}

	exists, err := images.Exists(*e.Container.clientPodman, e.Container.Image, &images.ExistsOptions{})
	if err != nil {
		return fmt.Errorf("failed to check if image exists: %v", err)
	}
	if !exists {
		if err := e.pullPodman(); err != nil {
			return err
		}
	}

	b := make([]byte, 4)
	rand.Read(b)
	shortid := hex.EncodeToString(b)

	// Ensure we can mount our filesystem to the same path inside the container
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir, err := filepath.Abs(wd)
	if err != nil {
		return err
	}
	linuxDir := toLinuxPath(dir)

	currPlatform := platforms.DefaultSpec()
	currPlatform.OS = "linux"

	log.Debugf("using image %s", e.Container.Image)

	s := specgen.NewSpecGenerator(e.Container.Image, false)
	s.Name = fmt.Sprintf("crie-%s-%s", filepath.Base(e.Name()), shortid)
	s.Entrypoint = []string{}
	s.Command = []string{"/bin/sh", "-c", "tail -f /dev/null"}
	s.WorkDir = linuxDir
	s.UserNS = specgen.Namespace{
		NSMode: "keep-id",
	}
	s.Mounts = []spec.Mount{
		{
			Type:        "bind",
			Source:      dir,
			Destination: linuxDir,
			Options:     []string{"rbind", "rw", "Z"},
		},
	}

	createResponse, err := containers.CreateWithSpec(*e.Container.clientPodman, s, nil)
	if err != nil {
		return err
	}
	e.Container.id = createResponse.ID

	startOptions := containers.StartOptions{}
	if err := containers.Start(*e.Container.clientPodman, createResponse.ID, &startOptions); err != nil {
		return err
	}

	return nil
}

func (e *Lint) pullPodman() error {
	_, err := images.Pull(*e.Container.clientPodman, e.Container.Image, nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *Lint) execPodman(params []string, stdout io.Writer) error {
	cmd := append([]string{e.Bin}, params...)

	log.Debug(cmd)
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	execCreateConfig := handlers.ExecCreateConfig{
		ExecOptions: container.ExecOptions{
			User:         currentUser.Uid,
			Cmd:          cmd,
			Privileged:   false,
			AttachStdin:  false,
			AttachStderr: true,
			AttachStdout: true,
			Tty:          false,
		},
	}
	execID, err := containers.ExecCreate(*e.Container.clientPodman, e.Container.id, &execCreateConfig)
	if err != nil {
		return err
	}

	logs, err := attachedExecStart(*e.Container.clientPodman, execID, &containers.ExecStartOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if _, err := io.Copy(stdout, logs); err != nil {
			log.Error("Error during reading logs: %v\n", err)
			return
		}
		err := logs.Close()
		if err != nil {
			log.Error(err)
		}
		_ = containers.ExecRemove(*e.Container.clientPodman, execID, &containers.ExecRemoveOptions{})
	}()

	timeout := time.After(5 * time.Second)
	check := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return errors.New("exec timed out")
		case <-check:
			inspect, err := containers.ExecInspect(*e.Container.clientPodman, execID, &containers.ExecInspectOptions{})
			if err != nil {
				return err
			}
			if inspect.Running == false {
				if inspect.ExitCode != 0 {
					return errors.New("exit code " + strconv.Itoa(inspect.ExitCode))
				}
				return nil
			}
		}
	}
}

func (e *Lint) cleanupPodman() error {

	// TODO cleanup based on labels (project, language)

	if e.Container.id != "" {
		var timeoutSeconds uint = 3
		var ignore = false

		d := log.WithFields(log.Fields{"podmanId": e.Container.id})

		d.Debug("stopping container")
		stopOptions := containers.StopOptions{
			Timeout: &timeoutSeconds,
			Ignore:  &ignore,
		}
		err := containers.Stop(*e.Container.clientPodman, e.Container.id, &stopOptions)
		if err != nil {
			return err
		}

		d.Debug("removing container")
		removeOptions := containers.RemoveOptions{
			Timeout: &timeoutSeconds,
		}
		containers.Remove(*e.Container.clientPodman, e.Container.id, &removeOptions)

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
		Detach bool `json:"Detach"`
	}{
		Detach: false,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := conn.DoRequest(ctx, bytes.NewReader(bodyJSON), http.MethodPost, "/exec/%s/start", nil, nil, sessionID)
	if resp == nil {
		return nil, err
	}

	return resp.Body, err
}
