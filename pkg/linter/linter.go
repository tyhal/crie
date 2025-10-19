package linter

import (
	"io"
	"sync"
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
	WillRun() error
	Cleanup(wg *sync.WaitGroup)
	MaxConcurrency() int
	Run(filePath string, rep chan Report)
}
