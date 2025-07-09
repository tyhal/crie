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
	"github.com/containers/podman/v5/pkg/machine/define"
	"github.com/containers/podman/v5/pkg/machine/env"
	"github.com/containers/podman/v5/pkg/machine/provider"
	"github.com/containers/podman/v5/pkg/machine/vmconfigs"
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

type podmanExecutor struct {
	Name   string
	Image  string
	client *context.Context
	id     string
}

var podmanInstalled = false

func willPodman() error {

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

func getPodmanMachineSocket() (socketPath string, err error) {
	currProvider, err := provider.Get()
	if err != nil {
		return
	}
	dirs, err := env.GetMachineDirs(currProvider.VMType())
	if err != nil {
		return
	}
	mc, vmErr := vmconfigs.LoadMachineByName(define.DefaultMachineName, dirs)
	if vmErr != nil {
		currentUser, userErr := user.Current()
		if userErr != nil {
			return socketPath, userErr
		}
		socketPath = fmt.Sprintf("/run/user/%s/podman/podman.sock", currentUser.Uid)
		return
	}
	podmanSocket, _, err := mc.ConnectionInfo(currProvider.VMType())
	if err != nil {
		return
	}
	socketPath = podmanSocket.Path
	return
}

func (e *podmanExecutor) setup() error {

	ctx := context.Background()

	{
		var uri string
		if _, found := os.LookupEnv("CONTAINER_HOST"); found {
			uri = ""
		} else {
			socketPath, err := getPodmanMachineSocket()
			if err != nil {
				return err
			}
			uri = fmt.Sprintf("unix://%s", socketPath)
		}

		c, err := bindings.NewConnection(ctx, uri)
		if err != nil {
			return err
		}
		e.client = &c
	}

	exists, err := images.Exists(*e.client, e.Image, &images.ExistsOptions{})
	if err != nil {
		return fmt.Errorf("failed to check if Image exists: %v", err)
	}
	if !exists {
		if err := e.pull(); err != nil {
			return err
		}
	}

	b := make([]byte, 4)
	rand.Read(b)
	shortid := hex.EncodeToString(b)

	wdHost, err := os.Getwd()
	if err != nil {
		return err
	}
	wdContainer, err := getwdContainer()
	if err != nil {
		return err
	}

	currPlatform := platforms.DefaultSpec()
	currPlatform.OS = "linux"

	log.Debugf("using image %s", e.Image)

	s := specgen.NewSpecGenerator(e.Image, false)
	s.Name = fmt.Sprintf("crie-%s-%s", filepath.Base(e.Name), shortid)
	s.Entrypoint = []string{"/bin/sh", "-c"}
	s.Command = []string{"tail -f /dev/null"}
	s.WorkDir = wdContainer
	s.UserNS = specgen.Namespace{
		NSMode: "keep-id",
	}
	s.Mounts = []spec.Mount{
		{
			Type:        "bind",
			Source:      wdHost,
			Destination: wdContainer,
			Options:     []string{"rbind", "rw", "Z"},
		},
	}

	createResponse, err := containers.CreateWithSpec(*e.client, s, nil)
	if err != nil {
		return err
	}
	e.id = createResponse.ID

	startOptions := containers.StartOptions{}
	if err := containers.Start(*e.client, e.id, &startOptions); err != nil {
		return err
	}

	return nil
}

func (e *podmanExecutor) pull() error {
	_, err := images.Pull(*e.client, e.Image, nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *podmanExecutor) exec(bin string, frontParams []string, filePath string, backParams []string, chdir bool, stdout io.Writer, _ io.Writer) error {

	wdContainer, err := getwdContainer()
	if err != nil {
		return err
	}
	if chdir {
		wdContainer = filepath.Join(wdContainer, filepath.Dir(filePath))
	}
	targetFile := filePath
	if chdir {
		targetFile = filepath.Base(filePath)
	}

	cmd := append([]string{bin}, frontParams...)
	cmd = append(cmd, targetFile)
	cmd = append(cmd, backParams...)

	log.Debug(cmd)
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	execCreateConfig := handlers.ExecCreateConfig{
		ExecOptions: container.ExecOptions{
			User:         currentUser.Uid,
			Cmd:          cmd,
			WorkingDir:   wdContainer,
			Privileged:   false,
			AttachStdin:  false,
			AttachStderr: true,
			AttachStdout: true,
			Tty:          false,
		},
	}

	execID, err := containers.ExecCreate(*e.client, e.id, &execCreateConfig)
	if err != nil {
		return err
	}

	logs, err := attachedExecStart(*e.client, execID, &containers.ExecStartOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if _, err := io.Copy(stdout, logs); err != nil {
			log.Errorf("Error during reading logs: %v\n", err)
			return
		}
		err := logs.Close()
		if err != nil {
			log.Error(err)
		}
		_ = containers.ExecRemove(*e.client, execID, &containers.ExecRemoveOptions{})
	}()

	timeout := time.After(5 * time.Second)
	check := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return errors.New("exec timed out")
		case <-check:
			inspect, err := containers.ExecInspect(*e.client, execID, &containers.ExecInspectOptions{})
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

func (e *podmanExecutor) cleanup() error {

	// TODO cleanup based on labels (project, language)

	if e.id != "" {
		var timeoutSeconds uint = 3
		var ignore = false

		d := log.WithFields(log.Fields{"podmanId": e.id})

		d.Debug("stopping container")
		stopOptions := containers.StopOptions{
			Timeout: &timeoutSeconds,
			Ignore:  &ignore,
		}
		err := containers.Stop(*e.client, e.id, &stopOptions)
		if err != nil {
			return err
		}

		d.Debug("removing container")
		removeOptions := containers.RemoveOptions{
			Timeout: &timeoutSeconds,
		}
		_, err = containers.Remove(*e.client, e.id, &removeOptions)
		if err != nil {
			return err
		}

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
