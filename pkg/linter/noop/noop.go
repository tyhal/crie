package noop

import (
	"math"
	"sync"

	"github.com/tyhal/crie/pkg/linter"
)

// LintNoop performs no operations, as a template implementation of cries' linter.Linter
type LintNoop struct {
	Type       string `json:"type" yaml:"type" jsonschema:"enum=noop" jsonschema_description:"a linter type to do nothing"`
	willRunErr error
	runErr     error
}

// WithErr creates a new LintNoop with the given errors
func WithErr(willRunErr error, runErr error) *LintNoop {
	return &LintNoop{
		Type:       "noop",
		runErr:     runErr,
		willRunErr: willRunErr,
	}
}

var _ linter.Linter = (*LintNoop)(nil)

// Name returns the name of the linter
func (l *LintNoop) Name() string {
	return "noop"
}

// WillRun just returns the configured error
func (l *LintNoop) WillRun() (err error) {
	return l.willRunErr
}

// Cleanup removes any additional resources created in the process
func (l *LintNoop) Cleanup(group *sync.WaitGroup) {
	group.Done()
}

// MaxConcurrency return max number of parallel files to fmt
func (l *LintNoop) MaxConcurrency() int {
	return math.MaxInt32
}

// Run will just return the configured error as a report
func (l *LintNoop) Run(filepath string, rep chan linter.Report) {
	rep <- linter.Report{File: filepath, Err: l.runErr, StdOut: nil, StdErr: nil}
}
