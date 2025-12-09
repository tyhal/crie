package orchestrator

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"runtime/trace"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/linter"
)

type Job struct {
	lint linter.Linter
	file string
}

type JobOrchestrator struct {
	Dispatchers  sync.WaitGroup
	executors    sync.WaitGroup
	cleanupStart sync.WaitGroup
	cleaners     sync.WaitGroup
	files        []string
	jobQ         chan Job
	report       sync.WaitGroup
	repQ         chan linter.Report
	reporter     linter.Reporter
	maxExecutors int
}

func New(files []string, reporter linter.Reporter) *JobOrchestrator {
	maxBacklog := 1024
	orch := &JobOrchestrator{
		files:        files,
		maxExecutors: min(runtime.NumCPU(), len(files)),
		jobQ:         make(chan Job, maxBacklog),
		repQ:         make(chan linter.Report, maxBacklog),
		reporter:     reporter,
	}
	orch.cleanupStart.Add(1)
	return orch
}

func (d *JobOrchestrator) executor(ctx context.Context) {
	for job := range d.jobQ {
		if job.lint == nil {
			log.Error("oh no")
		} else {
			// TODO lock on Format
			d.repQ <- job.lint.Run(job.file)
		}
	}
}

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
			d.executor(ctx)
		})
	}
	return d.wait
}

// Dispatcher submits jobQ to the workers
func (d *JobOrchestrator) Dispatcher(ctx context.Context, l linter.Linter, reg *regexp.Regexp) bool {
	defer trace.StartRegion(ctx, "A Dispatcher "+reg.String()).End()
	var startup sync.WaitGroup
	var active bool
	var matched []string
	for _, file := range d.files {
		if reg.MatchString(file) {
			if !active {
				active = true
				startup.Go(func() {
					err := l.Setup(ctx)
					if err != nil {
						log.Error(err)
						return
					}
					d.cleaners.Go(func() {
						d.cleanupStart.Wait()
						l.Cleanup(ctx)
						// TODO cleanup should probably call cancel on the context of the executors related to this linter
					})
				})
			}
			matched = append(matched, file)
		}
	}
	startup.Wait()
	if active {
		for _, file := range matched {
			d.jobQ <- Job{lint: l, file: file}
		}
	}
	return active
}

func (d *JobOrchestrator) wait() {
	d.Dispatchers.Wait()
	close(d.jobQ)
	d.executors.Wait()
	d.cleanupStart.Done()
	close(d.repQ)
	d.report.Wait()
	log.Info("waiting for cleanup")
	d.cleaners.Wait()
}
