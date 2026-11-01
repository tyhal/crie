// Package noop provides a linter that does nothing, useful for testing and as a template implementation of linter.Linter
package noop

import (
	"context"
	"errors"
	"runtime/trace"
	"sync/atomic"
	"time"

	"github.com/tyhal/crie/pkg/linter"
)

// ErrMissedSetup is returned when a linter wasn't setup before running or cleaning up
var ErrMissedSetup = errors.New("setup was not called first")

// LintNoop performs no operations, as a template implementation of cries' linter.Linter
type LintNoop struct {
	Type          string `json:"type" yaml:"type" jsonschema:"a linter type to do nothing"`
	setupErr      error
	runErr        error
	lintDuration  time.Duration
	setupDuration time.Duration
	setup         atomic.Bool
}

// WithErr creates a new LintNoop with the given errors
func WithErr(setupErr error, runErr error) *LintNoop {
	return &LintNoop{
		Type:     "noop",
		runErr:   runErr,
		setupErr: setupErr,
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

// Setup just returns the configured error
func (l *LintNoop) Setup(ctx context.Context) error {
	if l.setupDuration > 0 {
		defer trace.StartRegion(ctx, "Setup").End()
		time.Sleep(l.setupDuration)
	}
	l.setup.Store(true)
	return l.setupErr
}

// Cleanup removes any additional resources created in the process
func (l *LintNoop) Cleanup(ctx context.Context) error {
	if !l.setup.Load() {
		return ErrMissedSetup
	}
	if l.setupDuration > 0 {
		defer trace.StartRegion(ctx, "Cleanup").End()
		time.Sleep(l.setupDuration)
	}
	return nil
}

// Run will just return the configured error as a report
func (l *LintNoop) Run(ctx context.Context, filepath string) linter.Report {
	if !l.setup.Load() {
		return linter.Report{Target: filepath, File: filepath, Err: ErrMissedSetup, StdOut: nil, StdErr: nil}
	}
	if l.lintDuration > 0 {
		defer trace.StartRegion(ctx, "Run "+filepath).End()
		time.Sleep(l.lintDuration)
	}
	return linter.Report{Target: filepath, File: filepath, Err: l.runErr, StdOut: nil, StdErr: nil}
}
