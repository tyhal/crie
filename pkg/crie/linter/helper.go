package linter

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

// GetName getter method for any Linter's name
func GetName(l Linter) string {
	if l == nil {
		return ""
	}
	return l.Name()
}

func reporter(maxCon int, fileList []string, didRepErr chan bool, linterReport chan Report) {
	lintErr := false
	for i := range fileList {
		if i%maxCon == 0 {
			didRepErr <- lintErr
		}
		report := <-linterReport
		err := report.Log()
		if err != nil {
			log.WithFields(log.Fields{"type": "err"}).Debug(report.Err)
			lintErr = true
		}
	}
	didRepErr <- lintErr
}

// LintFileList simply takes a single Linter and runs it for each file
func LintFileList(l Linter, fileList []string) error {
	linterReport := make(chan Report)
	didRepErr := make(chan bool)
	maxCon := min(l.MaxConcurrency(), len(fileList))

	go reporter(maxCon, fileList, didRepErr, linterReport)

	didErr := false
	for i, codePath := range fileList {
		if i%maxCon == 0 {
			if <-didRepErr {
				didErr = true
			}
		}
		go l.Run(codePath, linterReport)
	}

	if <-didRepErr || didErr {
		return errors.New("some files failed to pass (" + l.Name() + ")")
	}
	return nil
}
