package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"runtime/trace"
	"sync"

	"github.com/tyhal/crie/pkg/errchain"
	"github.com/tyhal/crie/pkg/linter"
	"golang.org/x/sync/errgroup"
)

// Job is a single job to be run by the orchestrator
type Job struct {
	lint     linter.Linter
	file     string
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

// JobOrchestrator is responsible for dispatching jobs to workers
type JobOrchestrator struct {
	Dispatchers  sync.WaitGroup
	executors    sync.WaitGroup
	files        []string
	jobQ         chan Job
	report       errgroup.Group
	repQ         chan linter.Report
	reporter     linter.Reporter
	maxExecutors int
	locker       *fileLocker
	failFast     bool
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
		if job.lint == nil {
			d.repQ <- linter.Report{Err: fmt.Errorf("no linter for %s", job.file), Target: job.file}
		} else {
			d.repQ <- job.lint.Run(job.file)
			job.lintJobs.Done()
		}
	}
}

func (d *JobOrchestrator) lockExecutor() {
	for job := range d.jobQ {
		if job.lint == nil {
			d.repQ <- linter.Report{Err: errors.New("no linter found"), Target: job.file}
		} else {
			mu, ok := d.locker.locks[job.file]
			if !ok {
				d.repQ <- linter.Report{Err: errors.New("no lock found"), Target: job.file}
				continue
			}
			func() {
				mu.Lock()
				defer mu.Unlock()
				d.repQ <- job.lint.Run(job.file)
				job.lintJobs.Done()
			}()
		}
	}
}

// Start starts the orchestrator
func (d *JobOrchestrator) Start(ctx context.Context) func() error {
	d.report.Go(func() error {
		defer trace.StartRegion(ctx, "The Reporter").End()
		var anyErr error
		for report := range d.repQ {
			err := d.reporter.Log(&report)
			if err != nil {
				if d.failFast {
					return anyErr
				}
				if anyErr == nil {
					anyErr = errors.New("failures occurred")
				}
			}
		}
		return anyErr
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

	return d.wait
}

// Dispatcher submits jobQ to the workers
func (d *JobOrchestrator) Dispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp) {
	err := d.dispatcher(ctx, l, reg)
	if err != nil {
		d.repQ <- linter.Report{Err: err, Target: l.Name()}
	}
}

func (d *JobOrchestrator) dispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp) error {
	defer trace.StartRegion(ctx, "A Dispatcher "+reg.String()).End()

	var active bool
	var matched []string
	for _, file := range d.files {
		if reg.MatchString(file) {
			if !active {
				active = true
				err := l.Setup(ctx)
				if err != nil {
					return errchain.From(err).Link("failed to setup linter")
				}
			}
			matched = append(matched, file)
		}
	}
	if !active {
		return nil
	}

	var lintJobs sync.WaitGroup

	for _, file := range matched {
		lintJobs.Add(1)
		d.jobQ <- Job{lint: l, file: file, lintJobs: &lintJobs}
	}

	lintJobs.Wait()

	err := l.Cleanup(ctx)
	if err != nil {
		return errchain.From(err).Link("failed to cleanup linter")
	}

	return nil
}

func (d *JobOrchestrator) wait() error {
	d.Dispatchers.Wait()
	close(d.jobQ)
	d.executors.Wait()
	close(d.repQ)
	return d.report.Wait()
}
