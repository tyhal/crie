package exec

import (
	"context"
	"io"
)

// noopExecutor is a test/dummy executor that writes fixed strings to stdout and stderr.
type noopExecutor struct{}

// NewNoop creates an executor that does nothing
func NewNoop() Executor {
	return &noopExecutor{}
}

// Setup initializes the noopExecutor (no-op).
func (ne *noopExecutor) Setup(_ context.Context, _ Instance) error {
	return nil
}

// Exec writes sample output to stdout and stderr and returns nil.
func (ne *noopExecutor) Exec(_ string, stdout io.Writer, stderr io.Writer) error {
	_, err := stdout.Write([]byte("stdout"))
	if err != nil {
		return err
	}
	_, err = stderr.Write([]byte("stderr"))
	if err != nil {
		return err
	}
	return nil
}

// Cleanup finalizes the noopExecutor (no-op).
func (ne *noopExecutor) Cleanup(_ context.Context) error {
	return nil
}
