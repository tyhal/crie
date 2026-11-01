package executor

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/creack/pty"
	"golang.org/x/term"

	"github.com/tyhal/crie/pkg/linter"
)

// hostExecutor runs CLI tools directly on the host operating system.
type hostExecutor struct {
	Instance
	env []string
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
func (e *hostExecutor) Setup(_ context.Context, i Instance) error {
	e.Instance = i
	e.env = os.Environ()
	return nil
}

// command builds the exec.Cmd for a single run of the tool in dir
func (e *hostExecutor) command(ctx context.Context, dir string, targets []string) *exec.Cmd {
	c := exec.CommandContext(ctx, e.Bin, e.buildParams(targets...)...)
	c.Dir = dir
	// clipped so concurrent appends never share a backing array
	c.Env = slices.Clip(e.env)
	// appended entries win over any inherited TRACEPARENT from crie's own parent
	c.Env = append(c.Env, traceEnv(ctx)...)
	// exec only sets PWD implicitly when Env is nil, so mirror that here
	if pwd, err := filepath.Abs(dir); err == nil {
		c.Env = append(c.Env, "PWD="+pwd)
	}
	return c
}

// Exec runs the configured CLI tool on the host against the provided files.
func (e *hostExecutor) Exec(ctx context.Context, stdout io.Writer, stderr io.Writer, targets ...string) error {
	if len(targets) == 0 {
		return ErrNoTargets
	}
	dir, err := hostWorkdir(targets[0], e.ChDir, e.NoFileArg)
	if err != nil {
		return err
	}
	// batch output is split per file by its line prefixes, so it's kept plain rather than coloured by a tty
	// not for many files, a pty's colour and \r\n would get in the way of attributing their output
	if len(targets) == 1 && term.IsTerminal(int(os.Stdout.Fd())) {
		c := e.command(ctx, dir, targets)
		if ptm, err := pty.Start(c); err == nil {
			defer func() { _ = ptm.Close() }()
			_, _ = io.Copy(stderr, ptm)
			return linter.Result(c.Wait())
		}
	}

	c := e.command(ctx, dir, targets)
	c.Stdout = stdout
	c.Stderr = stderr
	return linter.Result(c.Run())
}

// Cleanup releases any resources allocated during host execution setup.
func (e *hostExecutor) Cleanup(_ context.Context) error {
	return nil
}
