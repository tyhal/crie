package linter

import (
	"io"
	"sync"
)

// FailedResultError means a linter didn't error but returned a failed result
type FailedResultError struct {
	err error
}

func (e *FailedResultError) Error() string {
	return e.Error()
}

// Result conditionally wraps an error with a FailedResultError or otherwise passes through nil, it should be used when a linter didn't error but returned a failed result
func Result(err error) error {
	if err == nil {
		return nil
	}
	return &FailedResultError{err: err}
}

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
	WillRun() error
	Cleanup(wg *sync.WaitGroup)
	MaxConcurrency() int
	Run(filePath string, rep chan Report)
}
