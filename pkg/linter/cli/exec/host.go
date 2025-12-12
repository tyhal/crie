package exec

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tyhal/crie/pkg/linter"
)

// HostExecutor runs CLI tools directly on the host operating system.
type HostExecutor struct {
	execCtx    context.Context
	execCancel context.CancelFunc
}

// WillHost checks whether the given binary can be found on the host PATH.
func WillHost(bin string) error {
	_, err := exec.LookPath(bin)
	return err
}

// Setup performs any required initialization for host execution.
func (e *HostExecutor) Setup(ctx context.Context) error {
	e.execCtx, e.execCancel = context.WithCancel(ctx)
	return nil
}

// Exec runs the configured CLI tool on the host against the provided file.
func (e *HostExecutor) Exec(i Instance, filePath string, stdout io.Writer, stderr io.Writer) error {
	targetFilePath := filePath
	if i.ChDir {
		targetFilePath = filepath.Base(filePath)
	}

	params := make([]string, 0, len(i.Start)+1+len(i.End))
	params = append(i.Start, targetFilePath)
	params = append(params, i.End...)

	c := exec.CommandContext(e.execCtx, i.Bin, params...)
	if i.ChDir {
		c.Dir = filepath.Dir(filePath)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.Dir = cwd
	}
	c.Stdout = stdout
	c.Stderr = stderr

	return linter.Result(c.Run())
}

// Cleanup releases any resources allocated during host execution setup.
func (e *HostExecutor) Cleanup(_ context.Context) error {
	if e.execCancel != nil {
		defer e.execCancel()
	}
	return nil
}
