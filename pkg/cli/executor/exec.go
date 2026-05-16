// Package executor allows for a CLI tool to be executed in different contexts
package executor

import (
	"context"
	"io"
	"path/filepath"
)

// Par represents command-line parameters passed to a linter executable.
type Par []string

// Instance describes a single CLI execution configuration for a linter tool.
// It specifies the binary to run, positional arguments to prepend/append,
// and whether execution should change to the file's directory first.
type Instance struct {
	Bin       string `json:"bin" yaml:"bin" jsonschema:"the binary or command to use"`
	Start     Par    `json:"start,omitzero" yaml:"start,flow" jsonschema:"parameters that will be put in front of the file path"`
	End       Par    `json:"end,omitzero" yaml:"end,flow" jsonschema:"parameters that will be put behind the file path"`
	ChDir     bool   `json:"chdir,omitzero" yaml:"chdir" jsonschema:"if true the tool will change directory to where the target file is located"`
	WillWrite bool   `json:"will_write,omitzero" yaml:"will_write" jsonschema:"whether the tool will write to the file system, if unset for container execution will mount with read-only privileges"`
	NoFileArg bool   `json:"no_file_arg,omitzero" yaml:"no_file_arg" jsonschema:"if true the target path is not inserted into command arguments; combine with chdir to use the target as the working directory"`
}

// buildParams constructs the command arguments after the binary name.
// targetPath should already be in the format expected by the executor (OS path for host, Linux path for containers).
func (i Instance) buildParams(targetPath string) []string {
	if i.NoFileArg {
		return append(append([]string{}, i.Start...), i.End...)
	}
	target := targetPath
	if i.ChDir {
		target = filepath.Base(targetPath)
	}
	return append(append(append([]string{}, i.Start...), target), i.End...)
}

// Executor abstracts how a CLI tool is executed (host, Docker, Podman, etc.).
// Implementations should prepare resources in Setup, execute the tool in Exec,
// and free resources in Cleanup.
type Executor interface {
	Setup(ctx context.Context, i Instance) error
	Exec(filepath string, stdout io.Writer, stderr io.Writer) error
	Cleanup(ctx context.Context) error
}
