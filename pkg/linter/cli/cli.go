package cli

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/crie/linter"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func toLinuxPath(dir string) string {
	splitPath := strings.Split(dir, ":")
	return filepath.ToSlash(splitPath[len(splitPath)-1])
}

// Name returns the command name
func (e *Lint) Name() string {
	return e.Bin
}

// WillRun does preflight checks for the 'Run'
func (e *Lint) WillRun() error {

	// Ensure cleanup channel is created
	if e.cleanedUp == nil {
		e.cleanedUp = make(chan error)
	}

	if e.execMode != auto {
		return errors.New("Crie doesn't support forcing a specific execution mode yet")
	}

	switch {
	case e.willPodman() == nil:
		e.execMode = podman
		if err := e.startPodman(); err != nil {
			return err
		}
	case e.willDocker() == nil:
		e.execMode = docker
		if err := e.startDocker(); err != nil {
			return err
		}
	case e.willHost() == nil:
		e.execMode = host
	default:
		return errors.New("could not determine execution mode")
	}

	log.Debugf("Using %s for %s", e.Name(), e.execMode)

	return nil
}

// Cleanup removes any additional resources created in the process
func (e *Lint) Cleanup() {
	var err error = nil
	defer func() {
		e.cleanedUp <- err
		close(e.cleanedUp)
	}()

	switch e.execMode {
	case podman:
		err = e.cleanupPodman()
	case docker:
		err = e.cleanupDocker()
	default:
	}
}

// WaitForCleanup Useful for when Cleanup is running in the background
func (e *Lint) WaitForCleanup() error {
	// TODO wait for cleanup should be the same for all linters

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

	switch e.execMode {
	case podman:
		err = e.execPodman(params, &outB)
	case docker:
		err = e.execDocker(params, &outB)
	case host:
		c := exec.Command(e.Bin, params...)
		c.Stdout = &outB
		c.Stderr = &errB
		err = c.Run()
	default:
		log.Error("Somehow we ran without determining the execution mode, doing nothing :(")
	}

	if e.execMode == podman || e.execMode == docker {
		rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: nil}
	} else {
		rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: &errB}
	}
}
