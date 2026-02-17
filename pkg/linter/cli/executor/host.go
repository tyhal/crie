package executor

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tyhal/crie/pkg/linter"
)

// hostExecutor runs CLI tools directly on the host operating system.
type hostExecutor struct {
	Instance
	execCtx    context.Context
	execCancel context.CancelFunc
}

// NewHost creates an executor that uses direct exec calls
func NewHost() Executor {
	return &hostExecutor{}
}

// WillHost checks whether the given binary can be found on the host PATH.
func WillHost(bin string) error {
	_, err := exec.LookPath(bin)
	return err
}

// Setup performs any required initialization for host execution.
func (e *hostExecutor) Setup(ctx context.Context, i Instance) error {
	e.Instance = i
	e.execCtx, e.execCancel = context.WithCancel(ctx)
	return nil
}

// Exec runs the configured CLI tool on the host against the provided file.
func (e *hostExecutor) Exec(filePath string, stdout io.Writer, stderr io.Writer) error {
	targetFilePath := filePath
	if e.ChDir {
		targetFilePath = filepath.Base(filePath)
	}

	params := make([]string, 0, len(e.Start)+1+len(e.End))
	params = append(e.Start, targetFilePath)
	params = append(params, e.End...)

	c := exec.CommandContext(e.execCtx, e.Bin, params...)
	if e.ChDir {
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
func (e *hostExecutor) Cleanup(_ context.Context) error {
	if e.execCancel != nil {
		defer e.execCancel()
	}
	return nil
}
