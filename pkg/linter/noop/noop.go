package noop

import (
	"math"
	"sync"

	"github.com/tyhal/crie/pkg/linter"
)

// LintNoop performs no operations, as a template implementation of cries' linter.Linter
type LintNoop struct {
	Type string `json:"type" yaml:"type" jsonschema:"enum=noop" jsonschema_description:"a linter type to do nothing"`
}

// Name returns the name of the linter
func (l *LintNoop) Name() string {
	return "noop"
}

// WillRun do nothing as there are no external deps
func (l *LintNoop) WillRun() (err error) {
	return nil
}

// Cleanup removes any additional resources created in the process
func (l *LintNoop) Cleanup(group *sync.WaitGroup) {
	group.Done()
}

// MaxConcurrency return max number of parallel files to fmt
func (l *LintNoop) MaxConcurrency() int {
	return math.MaxInt32
}

// Run the linter
func (l *LintNoop) Run(filepath string, rep chan linter.Report) {
	rep <- linter.Report{File: filepath, Err: nil, StdOut: nil, StdErr: nil}
}
