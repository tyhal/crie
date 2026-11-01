// Package linter provides a simple interface for linting, allowing for custom implementations
package linter

import (
	"context"
	"io"
)

// Report is used to state what issues a given file has
type Report struct {
	Target   string
	File     string // the file the report is about, when it's about one
	Err      error
	StdOut   io.Reader
	StdErr   io.Reader
	Findings []Finding
}

// Finding is a problem a tool reported at a position in a file, as read from output led by the file's path
type Finding struct {
	Tool    string
	File    string
	Line    int // 0 when the tool didn't give one
	Col     int // 0 when the tool didn't give one
	Message string
}

// Linter is a simple interface to enable a setup and check using WillRun before executing multiple Run's
// it is expected that Setup and Cleanup will be called exactly once given a context.
// Run may be called concurrently, its context is per file and carries the trace span for that run.
type Linter interface {
	Name() string
	Setup(ctx context.Context) error
	Cleanup(ctx context.Context) error
	Run(ctx context.Context, filePath string) Report
}

// BatchLinter is optionally implemented by a Linter that can lint many files in a single run.
// When Batched is true RunBatch is used instead of Run. It returns a Report for each file it could decide,
// and the files it couldn't (e.g. the batch failed but a file's own result is unknown) to be given their own Run,
// each with the findings the batch had for it to carry over to that Run's Report.
type BatchLinter interface {
	Linter
	Batched() bool
	RunBatch(ctx context.Context, filePaths []string) (reports []Report, rerun map[string][]Finding)
}
