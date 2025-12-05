// Package cli allows a CLI tool to be linted from the host or in a container
package cli

import (
	"bytes"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/linter"
	"github.com/tyhal/crie/pkg/linter/cli/exec"
)

// Version used to match img tags with crie versions
var Version = "latest"

// LintCli defines a predefined command to run against a file
type LintCli struct {
	Type           string        `json:"type" yaml:"type" jsonschema:"enum=cli" jsonschema_description:"the most common linter type, a cli tool"`
	Exec           exec.Instance `json:"exec" yaml:"exec" jsonschema_required:"true" jsonschema_description:"settings for the command to run" `
	Img            string        `json:"img,omitempty" yaml:"img,omitempty" jsonschema_description:"the container image to pull and use"`
	TagCrieVersion bool          `json:"tag_crie_version,omitempty" yaml:"tag_crie_version,omitempty" jsonschema_description:"if an image tag should be appended with cries current version"`
	executor       exec.Executor
}

var _ linter.Linter = (*LintCli)(nil)

func (e *LintCli) isContainer() bool {
	return e.Img != ""
}

// Name returns the command name
func (e *LintCli) Name() string {
	return e.Exec.Bin
}

func (e *LintCli) imgTagged() string {
	if e.TagCrieVersion {
		return fmt.Sprintf("%s:%s", e.Img, Version)
	}
	return e.Img
}

// WillRun does preflight checks for the 'Run'
func (e *LintCli) WillRun() error {

	img := e.imgTagged()

	switch {
	case e.isContainer() && exec.WillPodman() == nil:
		e.executor = &exec.PodmanExecutor{Name: e.Exec.Bin, Image: img}
	case e.isContainer() && exec.WillDocker() == nil:
		e.executor = &exec.DockerExecutor{Name: e.Exec.Bin, Image: img}
	case exec.WillHost(e.Exec.Bin) == nil:
		e.executor = &exec.HostExecutor{}
	default:
		return errors.New("could not determine execution mode [podman, docker, local]")
	}

	if err := e.executor.Setup(); err != nil {
		return err
	}

	return nil
}

// Cleanup removes any additional resources created in the process
func (e *LintCli) Cleanup() {
	if e.executor != nil {
		err := e.executor.Cleanup()
		if err != nil {
			log.Error(err)
		}
	}
}

// Run does the work required to lint the given filepath
func (e *LintCli) Run(filePath string) linter.Report {

	// Format any file received as an input.
	var outB, errB bytes.Buffer

	err := e.executor.Exec(e.Exec, filePath, &outB, &errB)

	return linter.Report{File: filePath, Err: err, StdOut: &outB, StdErr: &errB}
}
