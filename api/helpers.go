package api

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/api/linter"
	"io"
	"os"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func getName(l linter.Linter) string {
	if l == nil {
		return ""
	}
	return l.Name()
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// Narrows down the list by returning only results that do not match the match in the config file
func removeIgnored(list []string, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range list {
		result := f(entry)
		_, err := os.Stat(entry)
		if !result && err == nil {
			filteredLists = append(filteredLists, entry)
		}
	}
	return filteredLists
}

func filter(list []string, expect bool, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range list {
		result := f(entry)
		if result == expect {
			filteredLists = append(filteredLists, entry)
		}
	}
	return filteredLists
}

// LintFileList simply takes a single Linter and runs it for each file
func LintFileList(l linter.Linter, fileList []string) error {
	linterReport := make(chan linter.Report)
	didRepErr := make(chan bool)
	maxCon := min(maxConcurrency(), len(fileList))

	go func() {
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
	}()

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
