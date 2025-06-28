package cli

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tyhal/crie/pkg/crie/linter"
	"path/filepath"
	"strings"
	"time"
)

func (e *Lint) isContainer() bool {
	return e.Img != ""
}

// toLinuxPath ensures windows paths can be mapped to linux container paths
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

	// Ensure a cleanup channel is created
	if e.cleanedUp == nil {
		e.cleanedUp = make(chan error)
	}

	switch {
	case e.isContainer() && willPodman() == nil:
		e.executor = &podmanExecutor{Name: e.Bin, Image: e.Img}
	case e.isContainer() && willDocker() == nil:
		e.executor = &dockerExecutor{Name: e.Bin, Image: e.Img}
	case willHost(e.Bin) == nil:
		e.executor = &hostExecutor{}
	default:
		return errors.New("could not determine execution mode")
	}

	if err := e.executor.setup(); err != nil {
		return err
	}

	return nil
}

// Cleanup removes any additional resources created in the process
func (e *Lint) Cleanup() {
	var err error = nil
	defer func() {
		e.cleanedUp <- err
		close(e.cleanedUp)
	}()

	err = e.executor.cleanup()
}

// WaitForCleanup Useful for when Cleanup is running in the background
func (e *Lint) WaitForCleanup() error {
	var timeout time.Duration = 10

	select {
	case err := <-e.cleanedUp:
		return err
	case <-time.After(time.Second * timeout):
		return fmt.Errorf("timeout waiting for cleanup for %s (%d seconds)", e.Bin, timeout)
	}
}

// Run does the work required to lint the given filepath
func (e *Lint) Run(filePath string, rep chan linter.Report) {

	// Format any file received as an input.
	var outB, errB bytes.Buffer

	err := e.executor.exec(e.Bin, e.Start, toLinuxPath(filePath), e.End, e.ChDir, &outB, &errB)

	rep <- linter.Report{File: filePath, Err: err, StdOut: &outB, StdErr: &errB}
}
