package cli

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/containerd/platforms"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/crie/linter"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Lint defines a predefined command to run against a file
type Lint struct {
	Bin       string
	FrontPar  Par
	EndPar    Par
	Docker    DockerCmd
	useDocker bool
	cleanedUp chan error
}

// Par represents cli parameters
type Par []string

// DockerCmd an image to launch
type DockerCmd struct {
	Image  string
	client *client.Client
	id     string
}

// Name returns the command name
func (e *Lint) Name() string {
	return e.Bin
}

// WillRun does preflight checks for the 'Run'
func (e *Lint) WillRun() error {
	e.useDocker = exec.Command("which", e.Bin).Run() != nil

	// Ensure cleanup channel is created
	if e.cleanedUp == nil {
		e.cleanedUp = make(chan error)
	}

	// Startup sidecar container if needed
	if e.useDocker {
		if e.Docker.Image == "" {
			return errors.New("could not find " + e.Bin + ", possibly not installed")
		}
		if err := e.startDocker(); err != nil {
			return err
		}
	} else {
		log.Debug("using local binary")
	}

	return nil
}

// working solution posted to https://stackoverflow.com/questions/52145231/cannot-get-logs-from-docker-container-using-golang-docker-sdk
func (e *Lint) execDocker(params []string, stdout io.Writer) error {
	ctx := context.Background()
	cmd := append([]string{e.Bin}, params...)
	log.Trace(cmd)
	config := types.ExecConfig{
		Cmd:          cmd,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          false,
	}
	execResp, err := e.Docker.client.ContainerExecCreate(ctx, e.Docker.id, config)
	if err != nil {
		return err
	}

	startConfig := types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	}
	attach, err := e.Docker.client.ContainerExecAttach(ctx, execResp.ID, startConfig)
	if err != nil {
		return err
	}
	defer attach.Close()
	go func() {
		_, err := io.Copy(stdout, attach.Reader)
		if err != nil {
			log.Error(err)
		}
	}()

	timeout := time.After(5 * time.Second)
	check := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return errors.New("exec timed out")
		case <-check:
			inspect, err := e.Docker.client.ContainerExecInspect(ctx, execResp.ID)
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

func (e *Lint) pullDocker(ctx context.Context) error {

	// Ensure we have the image downloaded
	pullstat, err := e.Docker.client.ImagePull(ctx, e.Docker.Image, types.ImagePullOptions{})
	if err != nil {
		log.WithFields(log.Fields{
			"stage": "docker pull",
			"image": e.Docker.Image,
		}).Fatal(err)
		return err
	}

	var pullOut bytes.Buffer
	_, err = io.Copy(&pullOut, pullstat)
	if log.IsLevelEnabled(log.TraceLevel) {
		fmt.Print(pullOut.String())
	}
	return err
}

func toLinuxPath(dir string) string {
	splitPath := strings.Split(dir, ":")
	return filepath.ToSlash(splitPath[len(splitPath)-1])
}

func (e *Lint) startDocker() error {
	ctx := context.Background()

	// Add our client
	{
		c, err := client.NewClientWithOpts()
		if err != nil {
			return err
		}
		e.Docker.client = c
	}

	_, err := e.Docker.client.ImageHistory(ctx, e.Docker.Image)
	if err != nil {
		if err := e.pullDocker(ctx); err != nil {
			return err
		}
	}

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

	b := make([]byte, 4)
	rand.Read(b)
	shortid := hex.EncodeToString(b)

	resp, err := e.Docker.client.ContainerCreate(ctx,
		&container.Config{
			Entrypoint: []string{},
			Cmd:        []string{"/bin/sh", "-c", "tail -f /dev/null"},
			Image:      e.Docker.Image,
			WorkingDir: linuxDir,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: dir,
					Target: linuxDir,
				},
			},
		}, nil,
		&currPlatform,
		fmt.Sprintf("crie-%s-%s", filepath.Base(e.Name()), shortid))
	if err != nil {
		return err
	}
	e.Docker.id = resp.ID

	return e.Docker.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

// Cleanup removes any additional resources created in the process
func (e *Lint) Cleanup() {
	var err error = nil
	defer func() {
		e.cleanedUp <- err
		close(e.cleanedUp)
	}()

	if e.useDocker && e.Docker.id != "" {
		ctx := context.Background()
		timeoutSeconds := 3

		d := log.WithFields(log.Fields{"dockerId": e.Docker.id})

		d.Debug("stopping container")
		if err = e.Docker.client.ContainerStop(ctx, e.Docker.id, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
			return
		}
		d.Debug("removing container")
		if err = e.Docker.client.ContainerRemove(ctx, e.Docker.id, container.RemoveOptions{}); err != nil {
			return
		}
	}
}

// WaitForCleanup Useful for when Cleanup is running in the background
func (e *Lint) WaitForCleanup() error {
	var timeout time.Duration = 10

	select {
	case err := <-e.cleanedUp:
		return err
	case <-time.After(time.Second * timeout):
		return fmt.Errorf("timeout waiting for cleanup for %s (%d seconds)", e.Name(), timeout)
	}
}

// Run does the work required to lint the given filepath
func (e *Lint) Run(filepath string, rep chan linter.Report) {

	params := append(e.FrontPar, toLinuxPath(filepath))
	params = append(params, e.EndPar...)

	// Format any file received as an input.
	var outB, errB bytes.Buffer
	var err error
	if e.useDocker {
		err = e.execDocker(params, &outB)
	} else {
		// Local binary
		c := exec.Command(e.Bin, params...)
		c.Stdout = &outB
		c.Stderr = &errB
		err = c.Run()
	}

	rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: &errB}
}
