package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tyhal/crie/pkg/linter"
	"github.com/tyhal/crie/pkg/noop"
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
		expEndErr      string
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
			expCounter: reportCounter{0, 0}, // a failed setup aborts the run instead of reporting
			linter:     noop.WithErr(errors.New("setup err"), nil),
			dispatchers: []string{
				".*",
			},
			expEndErr: "setup err",
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
				ctx, end := orch.Start(t.Context())
				for _, regex := range tt.dispatchers {
					r, err := regexp.Compile(regex)
					assert.NoError(t, err)
					err = orch.CreateDispatcher(ctx, tt.linter, r, "")
					assert.ErrorIs(t, err, tt.expDispatchErr)
				}
				err := end()
				if tt.expEndErr != "" {
					assert.ErrorContains(t, err, tt.expEndErr)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, &tt.expCounter, counter)
			})
		})
	}
}

type failingReporter struct{ logged int }

func (r *failingReporter) Log(_ *linter.Report) error {
	r.logged++
	return errors.New("lint failed")
}

func TestFailFastSkipsQueuedJobs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		files := make([]string, 10*runtime.NumCPU())
		for i := range files {
			files[i] = fmt.Sprintf("f%d", i)
		}
		reporter := &failingReporter{}
		orch := New(files, reporter, false, true)
		ctx, end := orch.Start(t.Context())

		start := time.Now()
		assert.NoError(t, orch.CreateDispatcher(ctx, noop.WithSleep(time.Second, 0), regexp.MustCompile(".*"), ""))
		err := end()

		assert.ErrorContains(t, err, "lint failed")
		assert.Equal(t, 1, reporter.logged, "only the first failure should be reported")
		// executors may pick up one more batch before the failure is seen, but the rest of the queue is skipped
		assert.LessOrEqual(t, time.Since(start), 2*time.Second)
	})
}

// batchNoop is a linter.BatchLinter test double that records each batched and single run,
// files in rerun are left for their own run
type batchNoop struct {
	noop.LintNoop
	rerun      map[string]bool
	runs       atomic.Int32
	singleRuns atomic.Int32
}

var _ linter.BatchLinter = (*batchNoop)(nil)

func (b *batchNoop) Batched() bool { return true }

func (b *batchNoop) Run(ctx context.Context, filePath string) linter.Report {
	b.singleRuns.Add(1)
	return b.LintNoop.Run(ctx, filePath)
}

func (b *batchNoop) RunBatch(ctx context.Context, filePaths []string) (reps []linter.Report, rerun map[string][]linter.Finding) {
	b.runs.Add(1)
	for _, f := range filePaths {
		if b.rerun[f] {
			if rerun == nil {
				rerun = map[string][]linter.Finding{}
			}
			rerun[f] = []linter.Finding{{Tool: "noop", File: f, Line: 1, Message: "found in batch"}}
			continue
		}
		reps = append(reps, b.LintNoop.Run(ctx, f))
	}
	return reps, rerun
}

// targetRecorder is a reporter that records every reported target, and the findings reported for each
type targetRecorder struct {
	mu       sync.Mutex
	targets  []string
	findings map[string][]linter.Finding
}

func (r *targetRecorder) Log(rep *linter.Report) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.targets = append(r.targets, rep.Target)
	if len(rep.Findings) > 0 {
		if r.findings == nil {
			r.findings = map[string][]linter.Finding{}
		}
		r.findings[rep.Target] = rep.Findings
	}
	return rep.Err
}

func TestBatchDispatch(t *testing.T) {
	for _, locking := range []bool{false, true} {
		t.Run(fmt.Sprintf("locking=%v", locking), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				files := make([]string, 5*runtime.NumCPU()+3)
				for i := range files {
					files[i] = fmt.Sprintf("f%03d", i)
				}
				recorder := &targetRecorder{}
				orch := New(files, recorder, locking, false)
				ctx, end := orch.Start(t.Context())

				l := &batchNoop{LintNoop: noop.LintNoop{Type: "noop"}}
				assert.NoError(t, orch.CreateDispatcher(ctx, l, regexp.MustCompile(".*"), ""))
				assert.NoError(t, end())

				assert.ElementsMatch(t, files, recorder.targets, "every file should be reported once")
				assert.Equal(t, int32(1), l.runs.Load(), "files fitting in one batch should be linted in one run")
			})
		})
	}
}

func TestBatchReruns(t *testing.T) {
	for _, locking := range []bool{false, true} {
		t.Run(fmt.Sprintf("locking=%v", locking), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				files := []string{"a", "b", "c", "d", "e"}
				recorder := &targetRecorder{}
				orch := New(files, recorder, locking, false)
				ctx, end := orch.Start(t.Context())

				l := &batchNoop{LintNoop: noop.LintNoop{Type: "noop"}, rerun: map[string]bool{"b": true, "d": true}}
				assert.NoError(t, orch.CreateDispatcher(ctx, l, regexp.MustCompile(".*"), ""))
				assert.NoError(t, end())

				assert.ElementsMatch(t, files, recorder.targets, "every file should be reported once, reruns included")
				assert.Equal(t, int32(1), l.runs.Load())
				assert.Equal(t, int32(2), l.singleRuns.Load(), "each rerun file should get its own run")
				assert.Equal(t, map[string][]linter.Finding{
					"b": {{Tool: "noop", File: "b", Line: 1, Message: "found in batch"}},
					"d": {{Tool: "noop", File: "d", Line: 1, Message: "found in batch"}},
				}, recorder.findings, "a rerun's report should keep the findings its batch had for it")
			})
		})
	}
}

func TestBatchChunks(t *testing.T) {
	files := []string{"c/2", "a/1", "b/1", "a/2", "c/1"}
	sorted := slices.Sorted(slices.Values(files))

	// each path is 3 bytes plus a separator
	for maxBytes, want := range map[int]int{1 << 10: 1, 8: 3, 4: 5, 1: 5} {
		chunks := batchChunks(files, maxBytes)
		assert.Len(t, chunks, want, "maxBytes=%d", maxBytes)
		assert.Equal(t, sorted, slices.Concat(chunks...), "maxBytes=%d chunks should be contiguous and sorted", maxBytes)
		for _, chunk := range chunks {
			assert.NotEmpty(t, chunk)
		}
	}
	assert.Empty(t, batchChunks(nil, 4))
	assert.Equal(t, []string{"c/2", "a/1", "b/1", "a/2", "c/1"}, files, "input shouldn't be reordered")
}
