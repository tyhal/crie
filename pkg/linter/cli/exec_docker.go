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

var dockerInstalled = false

func (e *Lint) willDocker() error {
	if e.Container.Image == "" {
		return errors.New("no image specified for configuration " + e.Name() + "")
	}
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

func (e *Lint) startDocker() error {
	ctx := context.Background()

	// Add our clientDocker
	{
		c, err := client.NewClientWithOpts()
		if err != nil {
			return err
		}
		e.Container.clientDocker = c
	}

	_, err := e.Container.clientDocker.ImageHistory(ctx, e.Container.Image)
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

	resp, err := e.Container.clientDocker.ContainerCreate(ctx,
		&container.Config{
			Entrypoint: []string{},
			Cmd:        []string{"/bin/sh", "-c", "tail -f /dev/null"},
			Image:      e.Container.Image,
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
	e.Container.id = resp.ID

	return e.Container.clientDocker.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

func (e *Lint) pullDocker(ctx context.Context) error {

	// Ensure we have the image downloaded
	pullstat, err := e.Container.clientDocker.ImagePull(ctx, e.Container.Image, image.PullOptions{})
	if err != nil {
		log.WithFields(log.Fields{
			"stage": "docker pull",
			"image": e.Container.Image,
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

func (e *Lint) execDocker(params []string, stdout io.Writer) error {

	// working solution posted to https://stackoverflow.com/questions/52145231/cannot-get-logs-from-docker-container-using-golang-docker-sdk

	ctx := context.Background()
	cmd := append([]string{e.Bin}, params...)
	log.Trace(cmd)
	config := container.ExecOptions{
		Cmd:          cmd,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          false,
	}
	execResp, err := e.Container.clientDocker.ContainerExecCreate(ctx, e.Container.id, config)
	if err != nil {
		return err
	}

	startConfig := container.ExecAttachOptions{
		Detach: false,
		Tty:    false,
	}
	attach, err := e.Container.clientDocker.ContainerExecAttach(ctx, execResp.ID, startConfig)
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
			inspect, err := e.Container.clientDocker.ContainerExecInspect(ctx, execResp.ID)
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

func (e *Lint) cleanupDocker() error {

	// TODO cleanup based on labels (project, language)

	if e.Container.id != "" {
		ctx := context.Background()
		timeoutSeconds := 3

		d := log.WithFields(log.Fields{"dockerId": e.Container.id})

		d.Debug("stopping container")
		if err := e.Container.clientDocker.ContainerStop(ctx, e.Container.id, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
			return err
		}
		d.Debug("removing container")
		if err := e.Container.clientDocker.ContainerRemove(ctx, e.Container.id, container.RemoveOptions{}); err != nil {
			return err
		}
	}
	return nil
}
