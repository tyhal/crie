package noop

import (
	"math"
	"time"

	"github.com/tyhal/crie/pkg/linter"
)

// LintNoop performs no operations, as a template implementation of cries' linter.Linter
type LintNoop struct {
	Type          string `json:"type" yaml:"type" jsonschema:"enum=noop" jsonschema_description:"a linter type to do nothing"`
	willRunErr    error
	runErr        error
	lintDuration  time.Duration
	setupDuration time.Duration
}

// WithErr creates a new LintNoop with the given errors
func WithErr(willRunErr error, runErr error) *LintNoop {
	return &LintNoop{
		Type:       "noop",
		runErr:     runErr,
		willRunErr: willRunErr,
	}
}

// WithSleep creates a new LintNoop with the given sleep duration - useful for
func WithSleep(lint, setup time.Duration) *LintNoop {
	return &LintNoop{
		Type:          "noop",
		lintDuration:  lint,
		setupDuration: setup,
	}
}

var _ linter.Linter = (*LintNoop)(nil)

// Name returns the name of the linter
func (l *LintNoop) Name() string {
	return "noop"
}

// WillRun just returns the configured error
func (l *LintNoop) WillRun() (err error) {
	if l.setupDuration > 0 {
		time.Sleep(l.setupDuration)
	}
	return l.willRunErr
}

// Cleanup removes any additional resources created in the process
func (l *LintNoop) Cleanup() {
	if l.setupDuration > 0 {
		time.Sleep(l.setupDuration)
	}
}

// MaxConcurrency return max number of parallel files to fmt
func (l *LintNoop) MaxConcurrency() int {
	return math.MaxInt32
}

// Run will just return the configured error as a report
func (l *LintNoop) Run(filepath string) linter.Report {
	if l.lintDuration > 0 {
		time.Sleep(l.lintDuration)
	}
	return linter.Report{File: filepath, Err: l.runErr, StdOut: nil, StdErr: nil}
}
