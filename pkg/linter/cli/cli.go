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

// LintCli defines a predefined command to run against a file
type LintCli struct {
	Type     string `json:"type" yaml:"type" jsonschema:"enum=cli" jsonschema_description:"the most common linter type, a cli tool"`
	Bin      string `json:"bin" yaml:"bin" jsonschema_description:"the binary or command to use"`
	Start    Par    `json:"start,flow,omitempty" yaml:"start,flow,omitempty" jsonschema_description:"parameters that will be put in front of the file path"`
	End      Par    `json:"end,flow,omitempty" yaml:"end,flow,omitempty" jsonschema_description:"parameters that will be put behind the file path"`
	Img      string `json:"img,omitempty" yaml:"img,omitempty" jsonschema_description:"the container image to pull and use"`
	ChDir    bool   `json:"chdir,omitempty" yaml:"chdir,omitempty" jsonschema_description:"if true the tool will change directory to where the target file is located"`
	executor executor
}

// Par represents cli parameters
type Par []string

func (e *LintCli) isContainer() bool {
	return e.Img != ""
}

// toLinuxPath ensures windows paths can be mapped to linux container paths
func toLinuxPath(dir string) string {
	splitPath := strings.Split(dir, ":")
	return filepath.ToSlash(splitPath[len(splitPath)-1])
}

// Name returns the command name
func (e *LintCli) Name() string {
	return e.Bin
}

// WillRun does preflight checks for the 'Run'
func (e *LintCli) WillRun() error {

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
func (e *LintCli) Cleanup(group *sync.WaitGroup) {
	defer group.Done()
	err := e.executor.cleanup()
	if err != nil {
		log.Error(err)
	}
}

// Run does the work required to lint the given filepath
func (e *LintCli) Run(filePath string, rep chan linter.Report) {

	// Format any file received as an input.
	var outB, errB bytes.Buffer

	err := e.executor.exec(e.Bin, e.Start, toLinuxPath(filePath), e.End, e.ChDir, &outB, &errB)

	rep <- linter.Report{File: filePath, Err: err, StdOut: &outB, StdErr: &errB}
}
