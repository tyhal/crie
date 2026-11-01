package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"runtime"
	"runtime/trace"
	"slices"
	"sort"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"

	"github.com/tyhal/crie/pkg/linter"
)

var tracer = otel.Tracer("github.com/tyhal/crie/internal/runner/orchestrator")

// ErrBadDispatch is returned when a dispatcher is created with bad parameters
var ErrBadDispatch = errors.New("ctx, linter and regexp must be provided")

// Job is a single job to be run by the orchestrator
// it carries its dispatcher's context across the queue, like an http.Request, for the producer span and cancellation
type Job struct {
	ctx      context.Context
	lint     linter.Linter
	file     string
	files    []string         // when set, every file the run covers (a group_by module's, or a batch's) to lock instead of file
	findings []linter.Finding // what a batch found in the file this job reruns, kept over the rerun's own
	lintJobs *sync.WaitGroup
}

// fileLocker is used to lock files to prevent concurrent writes from formatters
type fileLocker struct {
	locks map[string]*sync.Mutex
}

// newFileLocker creates a new fileLocker from a file list
func newFileLocker(files []string) *fileLocker {
	m := make(map[string]*sync.Mutex, len(files))
	for _, f := range files {
		m[f] = &sync.Mutex{}
	}
	return &fileLocker{locks: m}
}

// JobOrchestrator is responsible for dispatching jobs to dispatchers
type JobOrchestrator struct {
	dispatchers  sync.WaitGroup
	executors    sync.WaitGroup
	files        []string
	jobQ         chan Job
	report       errgroup.Group
	repQ         chan linter.Report
	reporter     linter.Reporter
	maxExecutors int
	locker       *fileLocker
	failFast     bool
	// cancel aborts the run context returned by Start, with a dispatcher error or (when failFast) the first failed report
	cancel context.CancelCauseFunc
}

// New creates a new JobOrchestrator
// locking decides if we should have an exclusive lock per file due to potential writes
func New(files []string, reporter linter.Reporter, locking, failFast bool) *JobOrchestrator {
	maxBacklog := 1024
	var locker *fileLocker
	if locking {
		locker = newFileLocker(files)
	}
	orch := &JobOrchestrator{
		files:        files,
		maxExecutors: min(runtime.NumCPU(), len(files)),
		jobQ:         make(chan Job, maxBacklog),
		repQ:         make(chan linter.Report, maxBacklog),
		reporter:     reporter,
		locker:       locker,
		failFast:     failFast,
	}
	return orch
}

func (d *JobOrchestrator) executor() {
	for job := range d.jobQ {
		if job.ctx.Err() != nil {
			job.lintJobs.Done()
			continue
		}
		if job.lint == nil {
			d.repQ <- linter.Report{Err: fmt.Errorf("no linter for %s", job.file), Target: job.file}
		} else {
			d.finish(job, nil)
		}
	}
}

// finish runs job, reports its results and queues any of a batch's files that need their own run,
// then marks the job done
func (d *JobOrchestrator) finish(job Job, mutexes []*sync.Mutex) {
	defer job.lintJobs.Done()
	reps, rerun := d.run(job, mutexes)
	for _, rep := range reps {
		d.repQ <- rep
	}
	if len(rerun) == 0 {
		return
	}
	// queued from a goroutine as every executor blocking on a full jobQ would deadlock,
	// the dispatcher waits on lintJobs until they're all queued
	job.lintJobs.Go(func() {
		for _, file := range slices.Sorted(maps.Keys(rerun)) {
			if !d.enqueue(Job{ctx: job.ctx, lint: job.lint, file: file, findings: rerun[file], lintJobs: job.lintJobs}) {
				return
			}
		}
	})
}

// enqueue queues job unless its run is aborted first, reporting if it was queued
func (d *JobOrchestrator) enqueue(job Job) bool {
	job.lintJobs.Add(1)
	select {
	case d.jobQ <- job:
		return true
	case <-job.ctx.Done():
		job.lintJobs.Done()
		return false
	}
}

// run lints a single job under a consumer span, parented to the dispatcher's producer span,
// so everything the linter does for this target (container execs, child processes) is grouped beneath it.
// A batched job gives a report per file it could decide and the files to rerun on their own, otherwise there's just the one report.
func (d *JobOrchestrator) run(job Job, mutexes []*sync.Mutex) (reps []linter.Report, rerun map[string][]linter.Finding) {
	ctx, span := tracer.Start(job.ctx, "run "+job.lint.Name(),
		oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
		oteltrace.WithAttributes(attribute.String("crie.linter", job.lint.Name())),
	)
	defer span.End()
	if job.file != "" {
		span.SetAttributes(attribute.String("crie.target", job.file))
	}
	if len(job.files) > 0 {
		span.SetAttributes(attribute.Int("crie.files", len(job.files)))
	}

	if len(mutexes) > 0 {
		for _, mu := range mutexes {
			mu.Lock()
		}
		defer func() {
			for _, mu := range mutexes {
				mu.Unlock()
			}
		}()
		span.AddEvent("file locks acquired", oteltrace.WithAttributes(attribute.Int("crie.locked_files", len(mutexes))))
	}

	if bl, ok := job.lint.(linter.BatchLinter); ok && bl.Batched() && len(job.files) > 0 {
		reps, rerun = bl.RunBatch(ctx, job.files)
		if len(rerun) > 0 {
			span.AddEvent("rerunning files", oteltrace.WithAttributes(attribute.Int("crie.rerun", len(rerun))))
		}
	} else {
		rep := job.lint.Run(ctx, job.file)
		// a batch's machine readable findings are more reliable than whatever a rerun printed for people
		if len(job.findings) > 0 {
			rep.Findings = job.findings
		}
		reps = []linter.Report{rep}
	}
	for _, rep := range reps {
		if rep.Err != nil {
			span.SetStatus(codes.Error, rep.Err.Error())
			break
		}
	}
	return reps, rerun
}

func (d *JobOrchestrator) lockExecutor() {
	for job := range d.jobQ {
		if job.ctx.Err() != nil {
			job.lintJobs.Done()
			continue
		}
		if job.lint == nil {
			d.repQ <- linter.Report{Err: errors.New("no linter found"), Target: job.file}
			continue
		}

		filesToLock := job.files
		if len(filesToLock) == 0 {
			filesToLock = []string{job.file}
		}
		// sort to acquire locks in a consistent order and avoid deadlocks
		sort.Strings(filesToLock)

		var mutexes []*sync.Mutex
		valid := true
		for _, f := range filesToLock {
			mu, ok := d.locker.locks[f]
			if !ok {
				d.repQ <- linter.Report{Err: errors.New("no lock found"), Target: job.file}
				valid = false
				break
			}
			mutexes = append(mutexes, mu)
		}
		if !valid {
			job.lintJobs.Done()
			continue
		}

		d.finish(job, mutexes)
	}
}

// Start starts the orchestrator, like errgroup.WithContext the returned context is cancelled when the run is
// aborted and is the one to pass to CreateDispatcher. The returned func waits for the run to finish.
func (d *JobOrchestrator) Start(ctx context.Context) (context.Context, func() error) {
	ctx, d.cancel = context.WithCancelCause(ctx)

	d.report.Go(func() error {
		defer trace.StartRegion(ctx, "The Reporter").End()
		var allErrs []error
		for report := range d.repQ {
			// once aborted keep draining so executors never block, but stop reporting
			if ctx.Err() != nil {
				continue
			}
			err := d.reporter.Log(&report)
			if err != nil {
				err = fmt.Errorf("failed on %s: %w", report.Target, err)
				allErrs = append(allErrs, err)
				if d.failFast {
					d.cancel(err)
				}
			}
		}
		return errors.Join(allErrs...)
	})

	for i := range d.maxExecutors {
		d.executors.Go(func() {
			defer trace.StartRegion(ctx, fmt.Sprintf("Executor %d", i)).End()
			if d.locker != nil {
				d.lockExecutor()
			} else {
				d.executor()
			}
		})
	}

	return ctx, func() error { return d.wait(ctx) }
}

// CreateDispatcher submits jobQ to the dispatchers, ctx should be the one returned by Start so it stops when the run is aborted.
// If groupBy is non-empty, matched files are grouped by nearest ancestor directory
// containing the named marker file (e.g. "go.mod"), and the linter runs once per group.
// Otherwise a batched linter.BatchLinter runs once per chunk of the matched files, see batchChunks.
func (d *JobOrchestrator) CreateDispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp, groupBy string) error {
	if l == nil || reg == nil || ctx == nil {
		return ErrBadDispatch
	}
	d.dispatchers.Go(func() {
		// a linter that can't set up means the run can't do what was asked, so abort it
		if err := d.dispatcher(ctx, l, reg, groupBy); err != nil {
			d.cancel(err)
		}
	})
	return nil
}

func (d *JobOrchestrator) dispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp, groupBy string) (err error) {
	defer trace.StartRegion(ctx, "A CreateDispatcher "+reg.String()).End()

	var matched []string
	for _, file := range d.files {
		if reg.MatchString(file) {
			matched = append(matched, file)
		}
	}
	if len(matched) == 0 || ctx.Err() != nil {
		return nil
	}

	// jobs are enqueued with this context, so each run's consumer span is a child of this producer span
	ctx, span := tracer.Start(ctx, "dispatch "+l.Name(),
		oteltrace.WithSpanKind(oteltrace.SpanKindProducer),
		oteltrace.WithAttributes(
			attribute.String("crie.linter", l.Name()),
			attribute.String("crie.file_match", reg.String()),
			attribute.String("crie.group_by", groupBy),
			attribute.Int("crie.files", len(matched)),
		),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	// cleanup runs even after a failed setup or an aborted run so no containers are left behind
	defer func() {
		if cErr := l.Cleanup(context.WithoutCancel(ctx)); cErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to cleanup linter: %w", cErr))
		}
	}()

	if err := l.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup '%s': %w", l.Name(), err)
	}
	span.AddEvent("setup complete")

	var lintJobs sync.WaitGroup
	defer lintJobs.Wait()

	if groupBy != "" {
		groups := groupByModule(matched, groupBy)
		if len(groups) == 0 {
			return fmt.Errorf("no %s found above any matched files for %s", groupBy, l.Name())
		}
		for root, files := range groups {
			if !d.enqueue(Job{ctx: ctx, lint: l, file: root, files: files, lintJobs: &lintJobs}) {
				break
			}
		}
	} else if bl, ok := l.(linter.BatchLinter); ok && bl.Batched() {
		for _, chunk := range batchChunks(matched, maxBatchArgBytes) {
			if !d.enqueue(Job{ctx: ctx, lint: l, files: chunk, lintJobs: &lintJobs}) {
				break
			}
		}
	} else {
		for _, file := range matched {
			if !d.enqueue(Job{ctx: ctx, lint: l, file: file, lintJobs: &lintJobs}) {
				break
			}
		}
	}
	span.AddEvent("jobs dispatched")

	return nil
}

// maxBatchArgBytes keeps a batch's file arguments well under ARG_MAX (1MiB on macOS, ~2MiB on Linux)
const maxBatchArgBytes = 128 << 10

// batchChunks splits sorted files into as few chunks as fit within maxBytes of arguments each.
// It's deliberately not one chunk per executor, as each run of a tool pays its startup again
// (revive reloads type info per process, measured 1.1s as one run vs 1.9-2.2s as 8 parallel runs)
func batchChunks(files []string, maxBytes int) [][]string {
	var chunks [][]string
	var chunk []string
	size := 0
	for _, f := range slices.Sorted(slices.Values(files)) {
		if len(chunk) > 0 && size+len(f)+1 > maxBytes {
			chunks = append(chunks, chunk)
			chunk, size = nil, 0
		}
		chunk = append(chunk, f)
		size += len(f) + 1
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func (d *JobOrchestrator) wait(ctx context.Context) error {
	d.dispatchers.Wait()
	close(d.jobQ)
	d.executors.Wait()
	close(d.repQ)
	err := d.report.Wait()
	// the abort cause is returned unless it's the failed report already in err
	if cause := context.Cause(ctx); cause != nil && !errors.Is(err, cause) {
		err = errors.Join(err, cause)
	}
	d.cancel(nil)
	return err
}
