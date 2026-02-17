package orchestrator

import (
	"errors"
	"regexp"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter"
	"github.com/tyhal/crie/pkg/linter/noop"
)

type reportCounter struct {
	repCount int
	errCount int
}

func (r *reportCounter) Log(rep *linter.Report) error {
	r.repCount++
	if rep.Err != nil {
		r.errCount++
	}
	return nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		files          []string
		dispatchers    []string
		locking        bool
		linter         linter.Linter
		expCounter     reportCounter
		expDispatchErr error
	}{
		{
			name:  "empty files",
			files: []string{},
			dispatchers: []string{
				".*",
			},
			linter:     &noop.LintNoop{},
			expCounter: reportCounter{0, 0},
		},
		{
			name:       "missing dispatchers",
			files:      []string{"a", "b", "c"},
			linter:     &noop.LintNoop{},
			expCounter: reportCounter{0, 0},
		},
		{
			name:       "basic",
			files:      []string{"a", "b", "c"},
			expCounter: reportCounter{3, 0},
			linter:     &noop.LintNoop{},
			dispatchers: []string{
				".*",
			},
		},
		{
			name:   "basic with locking",
			files:  []string{"lock", "lock", "something"},
			linter: &noop.LintNoop{},
			dispatchers: []string{
				".*",
				".*",
			},
			locking:    true,
			expCounter: reportCounter{6, 0},
		},
		{
			name:       "nil linter",
			files:      []string{"a"},
			expCounter: reportCounter{0, 0},
			linter:     nil,
			dispatchers: []string{
				".*",
			},
			expDispatchErr: ErrBadDispatch,
		},
		{
			name:       "setup failure linter",
			files:      []string{"a", "b"},
			expCounter: reportCounter{1, 1}, // fail for each dispatcher
			linter:     noop.WithErr(errors.New("setup err"), nil),
			dispatchers: []string{
				".*",
			},
		},
		{
			name:       "run fail linter",
			files:      []string{"a", "b"},
			expCounter: reportCounter{2, 2}, // fail for each file
			linter:     noop.WithErr(nil, errors.New("run err")),
			dispatchers: []string{
				".*",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				counter := &reportCounter{}
				orch := New(tt.files, counter, tt.locking, false)
				end := orch.Start(t.Context())
				for _, regex := range tt.dispatchers {
					r, err := regexp.Compile(regex)
					assert.NoError(t, err)
					err = orch.CreateDispatcher(t.Context(), tt.linter, r)
					assert.ErrorIs(t, err, tt.expDispatchErr)
				}
				err := end()
				assert.NoError(t, err)
				assert.Equal(t, &tt.expCounter, counter)
			})
		})
	}
}
