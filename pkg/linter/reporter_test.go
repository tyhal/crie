package linter

import (
	"fmt"
	"io"
	"os"
	"testing"

	log "charm.land/log/v2"
	"github.com/stretchr/testify/assert"
)

// disableLogging changes the logger output and returns a restore function.
// The caller is expected to defer the returned function.
func disableLogging() func() {
	log.SetOutput(io.Discard)
	return func() {
		log.SetOutput(os.Stderr)
	}
}

func TestReporters(t *testing.T) {
	// TODO capture logs and check them
	defer disableLogging()()

	toolErr := fmt.Errorf("tool error")
	lintErr := &FailedResultError{
		err: fmt.Errorf("lint error"),
	}

	reporters := []Reporter{
		NewStandardReporter(true),
		NewStructuredReporter(true),
		NewStandardReporter(false),
		NewStructuredReporter(false),
	}
	tests := []struct {
		name   string
		err    error
		report *Report
	}{
		{"basic",
			toolErr,
			&Report{
				Err: toolErr,
			},
		},
		{
			"tool error",
			toolErr,
			&Report{
				Err: toolErr,
			},
		},
		{
			"lint error",
			lintErr,
			&Report{
				Err: lintErr,
			},
		},
	}

	for _, reporter := range reporters {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%T-%s", reporter, tt.name), func(t *testing.T) {
				err := reporter.Log(tt.report)
				assert.ErrorIs(t, err, tt.err)
			})
		}
	}
}
