package executor

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/creack/pty"
	"golang.org/x/term"

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
	dir := filepath.Dir(filePath)
	if !e.ChDir {
		var err error
		if dir, err = os.Getwd(); err != nil {
			return err
		}
	} else if e.NoFileArg {
		dir = filePath
	}

	// When running interactively, open a PTY so the subprocess detects a real
	// terminal and emits colour, while we still capture its output for reporting.
	// PTY merges stdout+stderr; captured into stderr so the reporter shows it.
	if term.IsTerminal(int(os.Stdout.Fd())) {
		c := exec.CommandContext(e.execCtx, e.Bin, e.buildParams(filePath)...)
		c.Dir = dir
		if ptm, err := pty.Start(c); err == nil {
			defer func() { _ = ptm.Close() }()
			_, _ = io.Copy(stderr, ptm)
			return linter.Result(c.Wait())
		}
	}

	c := exec.CommandContext(e.execCtx, e.Bin, e.buildParams(filePath)...)
	c.Dir = dir
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
