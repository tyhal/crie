package cli

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/crie/linter"
	"path/filepath"
	"strings"
	"sync"
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

	switch {
	case e.isContainer() && willPodman() == nil:
		e.executor = &podmanExecutor{Name: e.Bin, Image: e.Img}
	case e.isContainer() && willDocker() == nil:
		e.executor = &dockerExecutor{Name: e.Bin, Image: e.Img}
	case willHost(e.Bin) == nil:
		e.executor = &hostExecutor{}
	default:
		return errors.New("could not determine execution mode [podman, docker, local]")
	}

	if err := e.executor.setup(); err != nil {
		return err
	}

	return nil
}

// Cleanup removes any additional resources created in the process
func (e *Lint) Cleanup(group *sync.WaitGroup) {
	defer group.Done()
	err := e.executor.cleanup()
	if err != nil {
		log.Error(err)
	}
}

// Run does the work required to lint the given filepath
func (e *Lint) Run(filePath string, rep chan linter.Report) {

	// Format any file received as an input.
	var outB, errB bytes.Buffer

	err := e.executor.exec(e.Bin, e.Start, toLinuxPath(filePath), e.End, e.ChDir, &outB, &errB)

	rep <- linter.Report{File: filePath, Err: err, StdOut: &outB, StdErr: &errB}
}
