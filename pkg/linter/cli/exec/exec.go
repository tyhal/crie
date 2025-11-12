// Package exec allows for a CLI tool to be executed in different contexts
package exec

import "io"

// Par represents command-line parameters passed to a linter executable.
type Par []string

// Instance describes a single CLI execution configuration for a linter tool.
// It specifies the binary to run, positional arguments to prepend/append,
// and whether execution should change to the file's directory first.
type Instance struct {
	Bin   string `json:"bin" yaml:"bin" jsonschema_required:"true" jsonschema_description:"the binary or command to use"`
	Start Par    `json:"start,flow,omitempty" yaml:"start,flow,omitempty" jsonschema_description:"parameters that will be put in front of the file path"`
	End   Par    `json:"end,flow,omitempty" yaml:"end,flow,omitempty" jsonschema_description:"parameters that will be put behind the file path"`
	ChDir bool   `json:"chdir,omitempty" yaml:"chdir,omitempty" jsonschema_description:"if true the tool will change directory to where the target file is located"`
}

// Executor abstracts how a CLI tool is executed (host, Docker, Podman, etc.).
// Implementations should prepare resources in Setup, execute the tool in Exec,
// and free resources in Cleanup.
type Executor interface {
	Setup() error
	Exec(i Instance, filepath string, stdout io.Writer, stderr io.Writer) error
	Cleanup() error
}
