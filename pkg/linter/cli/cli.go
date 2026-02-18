// Package cli allows a CLI tool to be linted from the host or in a container
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime/trace"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/linter"
	"github.com/tyhal/crie/pkg/linter/cli/executor"
)

// Version used to match img tags with crie versions
var Version = "latest"

// LintCli defines a predefined command to run against a file
type LintCli struct {
	Type           string            `json:"type" yaml:"type" jsonschema:"enum=cli" jsonschema_description:"the most common linter type, a cli tool"`
	Exec           executor.Instance `json:"exec" yaml:"exec" jsonschema_required:"true" jsonschema_description:"settings for the command to run" `
	Img            string            `json:"img,omitempty" yaml:"img,omitempty" jsonschema_description:"the container image to pull and use"`
	TagCrieVersion bool              `json:"tag_crie_version,omitempty" yaml:"tag_crie_version,omitempty" jsonschema_description:"if an image tag should be appended with cries current version"`
	executor       executor.Executor
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

// Setup does preflight checks for the 'Run'
func (e *LintCli) Setup(ctx context.Context) error {
	defer trace.StartRegion(ctx, "Setup").End()

	img := e.imgTagged()

	switch {
	case e.isContainer() && executor.WillPodman(ctx) == nil:
		e.executor = executor.NewPodman(img)
	case e.isContainer() && executor.WillDocker() == nil:
		e.executor = executor.NewDocker(img)
	case executor.WillHost(e.Exec.Bin) == nil:
		e.executor = executor.NewHost()
	default:
		return errors.New("could not determine execution mode [podman, docker, local]")
	}

	if err := e.executor.Setup(ctx, e.Exec); err != nil {
		return fmt.Errorf("setting up executor %T: %w", e.executor, err)
	}

	return nil
}

// Cleanup removes any additional resources created in the process
func (e *LintCli) Cleanup(ctx context.Context) error {
	defer trace.StartRegion(ctx, "Cleanup").End()

	if e.executor != nil {
		err := e.executor.Cleanup(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

// Run does the work required to lint the given filepath
func (e *LintCli) Run(filePath string) linter.Report {
	// Format any file received as an input.
	var outB, errB bytes.Buffer

	err := e.executor.Exec(filePath, &outB, &errB)

	return linter.Report{Target: filePath, Err: err, StdOut: &outB, StdErr: &errB}
}
