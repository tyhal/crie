package exec

import (
	"io"
)

// Executor is an abstraction to allow any cli tool to run anywhere
type Executor interface {
	Setup() error
	Exec(bin string, frontParams []string, filePath string, endParams []string, chdir bool, stdout io.Writer, stderr io.Writer) error
	Cleanup() error
}
