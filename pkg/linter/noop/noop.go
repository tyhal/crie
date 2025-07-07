package noop

import (
	"github.com/tyhal/crie/pkg/crie/linter"
	"math"
	"sync"
)

type Lint struct{}

// Name returns the name of the linter
func (l *Lint) Name() string {
	return "noop"
}

// WillRun do nothing as there are no external deps
func (l *Lint) WillRun() (err error) {
	return nil
}

// Cleanup removes any additional resources created in the process
func (l *Lint) Cleanup(group *sync.WaitGroup) {
	group.Done()
}

// MaxConcurrency return max number of parallel files to fmt
func (l *Lint) MaxConcurrency() int {
	return math.MaxInt32
}

// Run shfmt -w
func (l *Lint) Run(filepath string, rep chan linter.Report) {
	rep <- linter.Report{File: filepath, Err: nil, StdOut: nil, StdErr: nil}
}
