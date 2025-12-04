// Package runner contains the logic to run the linters
package runner

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/errchain"
	"github.com/tyhal/crie/pkg/linter"
)

func (s *RunConfiguration) getRunningLanguages() (map[string]*Language, error) {
	currentLangs := s.Languages
	if s.Options.Only != "" {
		if lang, ok := s.Languages[s.Options.Only]; ok {
			currentLangs = map[string]*Language{s.Options.Only: lang}
		} else {
			return nil, fmt.Errorf("language %s not found", s.Options.Only)
		}
	}
	return currentLangs, nil
}

// LinterReady is used to communicate the linter that is ready to run
type LinterReady struct {
	lang   string
	linter linter.Linter
}

// FilematchReady is used to communicate a file list is ready to be processed
type FilematchReady struct {
	lang  string
	files []string
}

// WorkloadReady is a combination of the data needed to start dispatching jobs
type WorkloadReady struct {
	files  []string
	linter linter.Linter
}

type Job struct {
	linter linter.Linter
	file   string
}

// jobExecutor completes jobs and
func jobExecutor(jobs chan Job, reports chan linter.Report) {
	for job := range jobs {
		if job.linter == nil {
			log.Error("oh no")
		} else {
			reports <- job.linter.Run(job.file)
		}
	}
}

// jobSubmitter is a simple job dispatcher when both requirements are met it pushes jobs into a channel
func jobSubmitter(jobs chan Job, filesReady chan FilematchReady, lintReady chan LinterReady) {
	active := make(map[string]WorkloadReady)
	submit := func(lang string, workload WorkloadReady) {
		for _, file := range workload.files {
			jobs <- Job{linter: workload.linter, file: file}
		}
	}
	for {
		select {
		case fm, ok := <-filesReady:
			if !ok {
				filesReady = nil
			}
			if workload, ok := active[fm.lang]; ok {
				workload.files = fm.files
				submit(fm.lang, workload)
			} else {
				active[fm.lang] = WorkloadReady{files: fm.files}
			}
		case lr, ok := <-lintReady:
			if !ok {
				lintReady = nil
			}
			if workload, ok := active[lr.lang]; ok {
				workload.linter = lr.linter
				submit(lr.lang, workload)
			} else {
				active[lr.lang] = WorkloadReady{linter: lr.linter}
			}
		}
		if filesReady == nil && lintReady == nil {
			break
		}
	}
	close(jobs)
}

func (s *RunConfiguration) runLinters(lintType LintType, fileList []string) error {

	// every channel feeding up to workers matches worker count so that workers are priority
	maxWorkers := 16
	maxWorkloadBacklog := maxWorkers * 4

	currentLangs, err := s.getRunningLanguages()
	if err != nil {
		return err
	}

	filesReady := make(chan FilematchReady, maxWorkloadBacklog)
	lintReady := make(chan LinterReady, maxWorkloadBacklog)

	var langStartupWG sync.WaitGroup
	var langFileMatchingWG sync.WaitGroup
	for langName, lang := range currentLangs {
		langFileMatchingWG.Go(func() {
			currLint, err := lang.GetLinter(lintType)
			if err != nil {
				log.Error(err)
				return
			}
			if currLint == nil {
				return
			}
			var hasMatched bool
			reg := lang.FileMatch
			var matched []string
			for _, file := range fileList {
				if reg.MatchString(file) {
					if !hasMatched {
						hasMatched = true
						langStartupWG.Go(func() {
							err := currLint.WillRun()
							if err != nil {
								log.Error(err)
								return
							}
							lintReady <- LinterReady{
								linter: currLint,
								lang:   langName,
							}
						})
					}
					matched = append(matched, file)
				}
			}
			if hasMatched {
				filesReady <- FilematchReady{
					files: matched,
					lang:  langName,
				}
			}
		})
	}

	jobs := make(chan Job, maxWorkloadBacklog)
	reports := make(chan linter.Report, maxWorkloadBacklog*16)
	// start workers
	var workersWG sync.WaitGroup
	for range maxWorkers {
		workersWG.Go(func() {
			jobExecutor(jobs, reports)
		})
	}
	// submit jobs
	go jobSubmitter(jobs, filesReady, lintReady)

	var r linter.Reporter
	if s.Options.StrictLogging {
		r = linter.NewStructuredReporter(s.Options.Passes)
	} else {
		r = linter.NewStandardReporter(s.Options.Passes)
	}
	var reportWG sync.WaitGroup
	reportWG.Go(func() {
		for report := range reports {
			err := r.Log(&report)
			log.Error(err)
		}
	})

	langFileMatchingWG.Wait()
	langStartupWG.Wait()
	close(filesReady)
	close(lintReady)
	workersWG.Wait()
	close(reports)
	reportWG.Wait()

	return nil
}

// Run is the generic way to run everything based on the package configuration
func (s *RunConfiguration) Run(lintType LintType) error {
	fileList, err := s.getFileList()
	if err != nil {
		return errchain.From(err).Link("getting files")
	}
	err = s.runLinters(lintType, fileList)
	if err != nil {
		return err
	}
	log.Println("\u26c5  " + lintType.String() + "'ing passed")
	return nil
}
