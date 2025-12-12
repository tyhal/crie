package linter

import (
	"context"
	"io"
)

// Report is used to state what issues a given file has
type Report struct {
	File   string
	Err    error
	StdOut io.Reader
	StdErr io.Reader
}

// Linter is a simple interface to enable a setup and check using WillRun before executing multiple Run's
// it is expected that Setup and Cleanup will be called exactly once given a context.
// It is up to the implementation to pass a context to the Run method if needed and to cancel from the cleanup.
type Linter interface {
	Name() string
	Setup(ctx context.Context) error
	Cleanup(ctx context.Context) error
	Run(filePath string) Report
}
