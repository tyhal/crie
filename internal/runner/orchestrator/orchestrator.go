package orchestrator

import (
	"fmt"
	"regexp"
	"runtime"
	"sync"
	"time"

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
	maxWorkers   int
}

func New(files []string, reporter linter.Reporter) *JobOrchestrator {
	maxBacklog := 1024
	orch := &JobOrchestrator{
		files:      files,
		maxWorkers: min(runtime.NumCPU(), len(files)),
		jobQ:       make(chan Job, maxBacklog),
		repQ:       make(chan linter.Report, maxBacklog),
		reporter:   reporter,
	}
	orch.cleanupStart.Add(1)
	return orch
}

func (d *JobOrchestrator) executor() {
	for job := range d.jobQ {
		if job.lint == nil {
			log.Error("oh no")
		} else {
			// TODO lock on Format
			d.repQ <- job.lint.Run(job.file)
		}
	}
}

func (d *JobOrchestrator) Start() func() {
	d.report.Go(func() {
		for report := range d.repQ {
			err := d.reporter.Log(&report)
			if err != nil {
				log.Error(err)
			}
		}
	})
	for range d.maxWorkers {
		d.executors.Go(func() {
			d.executor()
		})
	}
	return d.wait
}

// Dispatcher submits jobQ to the workers
func (d *JobOrchestrator) Dispatcher(l linter.Linter, reg *regexp.Regexp) bool {
	var startup sync.WaitGroup
	var active bool
	var matched []string
	for _, file := range d.files {
		if reg.MatchString(file) {
			if !active {
				active = true
				startup.Go(func() {
					err := l.WillRun()
					time.Sleep(time.Second * 1)
					if err != nil {
						log.Error(err)
						return
					}
				})
				d.cleaners.Go(func() {
					d.cleanupStart.Wait()
					l.Cleanup()
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
	fmt.Println("waiting for cleanup")
	d.cleaners.Wait()
}
