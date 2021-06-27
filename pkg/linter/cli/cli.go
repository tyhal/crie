package cli

import (
	"bytes"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	linter "github.com/tyhal/crie/pkg/crie/linter"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// Lint defines a predefined command to run against a file
type Lint struct {
	Bin       string
	FrontPar  Par
	EndPar    Par
	Docker    DockerCmd
	useDocker bool
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
	if e.useDocker {
		if e.Docker.Image == "" {
			return errors.New("could not find " + e.Bin + ", possibly not installed")
		}
		if err := e.startDocker(); err != nil {
			return err
		}
		log.Warn("it's more efficient to have " + e.Bin + " installed locally")
	}
	return nil
}

// working solution posted to https://stackoverflow.com/questions/52145231/cannot-get-logs-from-docker-container-using-golang-docker-sdk
func (e *Lint) execDocker(params []string, stdout io.Writer) error {
	ctx := context.Background()
	cmd := append([]string{"/bin/" + e.Bin}, params...)
	config := types.ExecConfig{
		Cmd:          cmd,
		Env:          os.Environ(),
		AttachStderr: true,
		AttachStdout: true,
		Tty:          false,
	}
	execResp, err := e.Docker.client.ContainerExecCreate(ctx, e.Docker.id, config)
	if err != nil {
		return err
	}

	attach, err := e.Docker.client.ContainerExecAttach(ctx, execResp.ID, config)
	if err != nil {
		return err
	}
	defer attach.Close()
	go io.Copy(stdout, attach.Reader)

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

func (e *Lint) startDocker() error {
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	// Ensure we have the image downloaded
	pullstat, err := c.ImagePull(ctx, e.Docker.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	var pullOut bytes.Buffer
	if _, err = io.Copy(&pullOut, pullstat); err != nil {
		panic(err)
	}
	log.Debug(pullOut.String())

	e.Docker.client = c
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	resp, err := e.Docker.client.ContainerCreate(ctx,
		&container.Config{
			Entrypoint: []string{"/bin/sleep"},
			Cmd:        []string{"3600"},
			Image:      e.Docker.Image,
			WorkingDir: dir,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: dir,
					Target: dir,
				},
			},
		}, nil,
		"crie-"+e.Name())
	if err != nil {
		return err
	}
	if err := e.Docker.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	e.Docker.id = resp.ID
	return nil
}

// DidRun should be called after all other Runs to clean up
func (e *Lint) DidRun() {
	if e.useDocker {
		ctx := context.Background()
		var timeout time.Duration
		timeout = time.Second * 3

		d := log.WithFields(log.Fields{"dockerId": e.Docker.id})

		d.Debug("stopping container")
		if err := e.Docker.client.ContainerStop(ctx, e.Docker.id, &timeout); err != nil {
			log.Error(err)
		}
		d.Debug("removing container")
		if err := e.Docker.client.ContainerRemove(ctx, e.Docker.id, types.ContainerRemoveOptions{}); err != nil {
			log.Fatal(err)
		}
	}
}

// Run does the work required to lint the given filepath
func (e *Lint) Run(filepath string, rep chan linter.Report) {

	params := append(e.FrontPar, filepath)
	params = append(params, e.EndPar...)

	// Format any file received as input.
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
