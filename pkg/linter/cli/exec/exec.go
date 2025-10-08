package exec

import "io"

// Par represents cli parameters
type Par []string

type ExecInstance struct {
	Bin   string `json:"bin" yaml:"bin" jsonschema_required:"true" jsonschema_description:"the binary or command to use"`
	Start Par    `json:"start,flow,omitempty" yaml:"start,flow,omitempty" jsonschema_description:"parameters that will be put in front of the file path"`
	End   Par    `json:"end,flow,omitempty" yaml:"end,flow,omitempty" jsonschema_description:"parameters that will be put behind the file path"`
	ChDir bool   `json:"chdir,omitempty" yaml:"chdir,omitempty" jsonschema_description:"if true the tool will change directory to where the target file is located"`
}

// Executor is an abstraction to allow any cli tool to run anywhere
type Executor interface {
	Setup() error
	Exec(i ExecInstance, filepath string, stdout io.Writer, stderr io.Writer) error
	Cleanup() error
}
