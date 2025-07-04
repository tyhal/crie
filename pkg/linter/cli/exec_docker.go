package cli

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/containerd/platforms"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type dockerExecutor struct {
	Name   string
	Image  string
	client *client.Client
	id     string
}

var dockerInstalled = false

func willDocker() error {
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

func (e *dockerExecutor) setup() error {
	ctx := context.Background()

	// Add our clientDocker
	{
		c, err := client.NewClientWithOpts()
		if err != nil {
			return err
		}
		e.client = c
	}

	_, err := e.client.ImageHistory(ctx, e.Image)
	if err != nil {
		if err := e.pull(ctx); err != nil {
			return err
		}
	}

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

	b := make([]byte, 4)
	rand.Read(b)
	shortid := hex.EncodeToString(b)

	resp, err := e.client.ContainerCreate(ctx,
		&container.Config{
			Entrypoint: []string{},
			Cmd:        []string{"/bin/sh", "-c", "tail -f /dev/null"},
			Image:      e.Image,
			WorkingDir: wdContainer,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: wdHost,
					Target: wdContainer,
				},
			},
		}, nil,
		&currPlatform,
		fmt.Sprintf("crie-%s-%s", filepath.Base(e.Name), shortid))
	if err != nil {
		return err
	}
	e.id = resp.ID

	return e.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

func (e *dockerExecutor) pull(ctx context.Context) error {

	// Ensure we have the Image downloaded
	pullstat, err := e.client.ImagePull(ctx, e.Image, image.PullOptions{})
	if err != nil {
		log.WithFields(log.Fields{
			"stage": "docker pull",
			"Image": e.Image,
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

func (e *dockerExecutor) exec(bin string, frontParams []string, filePath string, backParams []string, chdir bool, stdout io.Writer, _ io.Writer) error {

	// working solution posted to https://stackoverflow.com/questions/52145231/cannot-get-logs-from-docker-container-using-golang-docker-sdk

	ctx := context.Background()

	// Ensure we can mount our filesystem to the same path inside the container
	containerWD, err := getwdContainer()
	if err != nil {
		return err
	}
	if chdir {
		containerWD = filepath.Join(containerWD, filepath.Dir(filePath))
	}
	targetFile := filePath
	if chdir {
		targetFile = filepath.Base(filePath)
	}

	cmd := append([]string{bin}, frontParams...)
	cmd = append(cmd, targetFile)
	cmd = append(cmd, backParams...)

	log.Trace(cmd)
	config := container.ExecOptions{
		Cmd:          cmd,
		WorkingDir:   containerWD,
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
			inspect, err := e.client.ContainerExecInspect(ctx, execResp.ID)
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

func (e *dockerExecutor) cleanup() error {

	// TODO cleanup based on labels (project, language)

	if e.id != "" {
		ctx := context.Background()
		timeoutSeconds := 3

		d := log.WithFields(log.Fields{"dockerId": e.id})

		d.Debug("stopping container")
		if err := e.client.ContainerStop(ctx, e.id, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
			return err
		}
		d.Debug("removing container")
		if err := e.client.ContainerRemove(ctx, e.id, container.RemoveOptions{}); err != nil {
			return err
		}
	}
	return nil
}
