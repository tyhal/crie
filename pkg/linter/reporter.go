package linter

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/folding"
)

// Runner will handle parallel runs of linters
type Runner struct {
	ShowPass      bool
	StrictLogging bool
	folder        folding.Folder
}

func (r *Runner) logConditional(reader io.Reader, typeField string, level log.Level) {
	field := log.WithFields(log.Fields{"type": typeField})
	if r.StrictLogging || level == log.DebugLevel {
		field.Log(level, reader)
	} else {
		field.Log(level)
		fmt.Println(reader)
	}
}

// Log simple takes all fields and pushes them to our using the default logger
func (r *Runner) Log(rep *Report) error {
	if rep.Err == nil {
		if r.ShowPass {
			if r.StrictLogging {
				log.Printf("pass %v", rep.File)
			} else {
				fmt.Printf("\u2714 %v\n", rep.File)
			}
			r.logConditional(rep.StdOut, "stdout", log.DebugLevel)
		}
	} else {
		var id string
		if r.StrictLogging {
			log.Printf("fail %v", rep.File)
		} else {
			id, _ = r.folder.Start(rep.File, "\u2716", false)
		}
		var failedResultErr *FailedResultError
		if errors.As(rep.Err, &failedResultErr) {
			r.logConditional(rep.StdErr, "stderr", log.ErrorLevel)
			r.logConditional(rep.StdOut, "stdout", log.InfoLevel)
			r.logConditional(strings.NewReader(rep.Err.Error()), "toolerr", log.DebugLevel)
		} else {
			r.logConditional(strings.NewReader(rep.Err.Error()), "toolerr", log.ErrorLevel)
		}
		if !r.StrictLogging {
			_ = r.folder.Stop(id)
		}
	}
	return rep.Err
}

func (r *Runner) listen(results chan error, linterReport chan Report) {
	r.folder = folding.New()
	for report := range linterReport {
		err := r.Log(&report)
		results <- err
	}
	close(results)
}

// LintFileList simply takes a single Linter and runs it for each file
func (r *Runner) LintFileList(l Linter, fileList []string) error {
	maxCon := min(l.MaxConcurrency(), len(fileList))
	if maxCon <= 0 {
		maxCon = 1
	}

	// job queue and response
	jobs := make(chan string, len(fileList))

	// reporting routine
	linterReport := make(chan Report, maxCon)
	results := make(chan error, len(fileList))
	go r.listen(results, linterReport)

	// create workers
	var wg sync.WaitGroup
	for i := 0; i < maxCon; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				path, ok := <-jobs
				if !ok {
					return // Channel closed, exit worker
				}
				l.Run(path, linterReport)
			}
		}()
	}

	// submit jobs
	for _, codePath := range fileList {
		jobs <- codePath
	}
	close(jobs)
	// wait for workers to exit
	go func() {
		wg.Wait()
		close(linterReport)
	}()

	// read results
	hasErrors := false
	for i := 0; i < len(fileList); i++ {
		if err := <-results; err != nil {
			hasErrors = true
		}
	}

	if hasErrors {
		return errors.New("some files failed to pass (" + l.Name() + ")")
	}
	return nil
}
