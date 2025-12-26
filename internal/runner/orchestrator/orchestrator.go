package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"runtime/trace"
	"sync"

	log "github.com/sirupsen/logrus"
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

type FileLocker struct {
	locks map[string]*sync.Mutex
}

func NewFileLocker(files []string) *FileLocker {
	m := make(map[string]*sync.Mutex, len(files))
	for _, f := range files {
		m[f] = &sync.Mutex{}
	}
	return &FileLocker{locks: m}
}

// JobOrchestrator is responsible for dispatching jobs to workers
type JobOrchestrator struct {
	Dispatchers  sync.WaitGroup
	executors    sync.WaitGroup
	files        []string
	jobQ         chan Job
	report       sync.WaitGroup
	repQ         chan linter.Report
	reporter     linter.Reporter
	maxExecutors int
	locker       *FileLocker
}

// New creates a new JobOrchestrator
// locking decides if we should have an exclusive lock per file due to potential writes
func New(files []string, reporter linter.Reporter, locking bool) *JobOrchestrator {
	maxBacklog := 1024
	var locker *FileLocker
	if locking {
		locker = NewFileLocker(files)
	}
	orch := &JobOrchestrator{
		files:        files,
		maxExecutors: min(runtime.NumCPU(), len(files)),
		jobQ:         make(chan Job, maxBacklog),
		repQ:         make(chan linter.Report, maxBacklog),
		reporter:     reporter,
		locker:       locker,
	}
	return orch
}

func (d *JobOrchestrator) executor() {
	for job := range d.jobQ {
		if job.lint == nil {
			d.repQ <- linter.Report{Err: fmt.Errorf("no linter for %s", job.file), File: job.file}
		} else {
			d.repQ <- job.lint.Run(job.file)
			job.lintJobs.Done()
		}
	}
}

func (d *JobOrchestrator) lockExecutor() {
	for job := range d.jobQ {
		if job.lint == nil {
			d.repQ <- linter.Report{Err: errors.New("no linter found"), File: job.file}
		} else {
			mu, ok := d.locker.locks[job.file]
			if !ok {
				d.repQ <- linter.Report{Err: errors.New("no lock for found"), File: job.file}
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
func (d *JobOrchestrator) Start(ctx context.Context) func() {
	d.report.Go(func() {
		defer trace.StartRegion(ctx, "The Reporter").End()
		for report := range d.repQ {
			err := d.reporter.Log(&report)
			if err != nil {
				log.Error(err)
			}
		}
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

	return d.Wait
}

// Dispatcher submits jobQ to the workers
func (d *JobOrchestrator) Dispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp) error {
	defer trace.StartRegion(ctx, "A Dispatcher "+reg.String()).End()

	// startup signals when the linter is ready to accept jobs
	var startup errgroup.Group

	var active bool
	var matched []string
	for _, file := range d.files {
		if reg.MatchString(file) {
			if !active {
				active = true
				startup.Go(func() error {
					return l.Setup(ctx)
				})
			}
			matched = append(matched, file)
		}
	}
	if !active {
		return nil
	}

	err := startup.Wait()
	if err != nil {
		return errchain.From(err).LinkF("failed to setup linter %s", l.Name())
	}

	var lintJobs sync.WaitGroup

	for _, file := range matched {
		lintJobs.Add(1)
		d.jobQ <- Job{lint: l, file: file, lintJobs: &lintJobs}
	}

	lintJobs.Wait()
	err = l.Cleanup(ctx)
	if err != nil {
		return errchain.From(err).LinkF("failed to cleanup linter %s", l.Name())
	}
	return nil
}

func (d *JobOrchestrator) Wait() {
	d.Dispatchers.Wait()
	close(d.jobQ)
	d.executors.Wait()
	close(d.repQ)
	d.report.Wait()
}
