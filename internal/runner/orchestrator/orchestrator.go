package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"runtime/trace"
	"sort"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/tyhal/crie/pkg/linter"
)

// ErrBadDispatch is returned when a dispatcher is created with bad parameters
var ErrBadDispatch = errors.New("ctx, linter and regexp must be provided")

// Job is a single job to be run by the orchestrator
type Job struct {
	lint      linter.Linter
	file      string
	lockFiles []string // individual files to lock; when set, used instead of file (for grouped runs)
	lintJobs  *sync.WaitGroup
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
			continue
		}

		filesToLock := job.lockFiles
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
			continue
		}

		func() {
			for _, mu := range mutexes {
				mu.Lock()
			}
			defer func() {
				for _, mu := range mutexes {
					mu.Unlock()
				}
			}()
			d.repQ <- job.lint.Run(job.file)
			job.lintJobs.Done()
		}()
	}
}

// Start starts the orchestrator
func (d *JobOrchestrator) Start(ctx context.Context) func() error {
	d.report.Go(func() error {
		defer trace.StartRegion(ctx, "The Reporter").End()
		var allErrs []error
		for report := range d.repQ {
			err := d.reporter.Log(&report)
			if err != nil {
				err = fmt.Errorf("failed on %s: %w", report.Target, err)
				if d.failFast {
					return err
				}
				allErrs = append(allErrs, err)
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

	return d.wait
}

// CreateDispatcher submits jobQ to the dispatchers.
// If groupBy is non-empty, matched files are grouped by nearest ancestor directory
// containing the named marker file (e.g. "go.mod"), and the linter runs once per group.
func (d *JobOrchestrator) CreateDispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp, groupBy string) error {
	if l == nil || reg == nil || ctx == nil {
		return ErrBadDispatch
	}
	d.dispatchers.Go(func() {
		err := d.dispatcher(ctx, l, reg, groupBy)
		if err != nil {
			d.repQ <- linter.Report{Err: err, Target: l.Name()}
		}
	})
	return nil
}

func (d *JobOrchestrator) dispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp, groupBy string) error {
	defer trace.StartRegion(ctx, "A CreateDispatcher "+reg.String()).End()

	var matched []string
	for _, file := range d.files {
		if reg.MatchString(file) {
			matched = append(matched, file)
		}
	}
	if len(matched) == 0 {
		return nil
	}

	if err := l.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup '%s': %w", l.Name(), err)
	}

	var lintJobs sync.WaitGroup

	if groupBy != "" {
		groups := groupByModule(matched, groupBy)
		if len(groups) == 0 {
			return fmt.Errorf("no %s found above any matched files for %s", groupBy, l.Name())
		}
		for root, files := range groups {
			lintJobs.Add(1)
			d.jobQ <- Job{lint: l, file: root, lockFiles: files, lintJobs: &lintJobs}
		}
	} else {
		for _, file := range matched {
			lintJobs.Add(1)
			d.jobQ <- Job{lint: l, file: file, lintJobs: &lintJobs}
		}
	}

	lintJobs.Wait()

	if err := l.Cleanup(ctx); err != nil {
		return fmt.Errorf("failed to cleanup linter: %w", err)
	}

	return nil
}

func (d *JobOrchestrator) wait() error {
	d.dispatchers.Wait()
	close(d.jobQ)
	d.executors.Wait()
	close(d.repQ)
	return d.report.Wait()
}
