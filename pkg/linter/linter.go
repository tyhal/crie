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
type Linter interface {
	Name() string
	WillRun(ctx context.Context) error
	Cleanup(ctx context.Context)
	Run(ctx context.Context, filePath string) Report
}
