package api

// TODO(tyler) this has repeated exec calls

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/api/linter"
)

// ContinueOnError Make flag
var ContinueOnError = false

// Verbose to print Report regardless of return code
var Verbose = false

// Quiet silence extra output
var Quiet = false

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

func printReportErr(rep linter.Report) error {
	if rep.Err == nil {
		fmt.Println(" ✔️  " + rep.File)
		return nil
	}

	fmt.Println(" ✖️  " + rep.File)

	if Quiet {
		return rep.Err
	}

	log.WithFields(log.Fields{"type": "stdout"}).Info(rep.StdOut)
	log.WithFields(log.Fields{"type": "stderr"}).Error(rep.StdErr)

	if ContinueOnError {
		log.Error(rep.Err.Error())
	} else {
		log.Fatal(rep.Err)
	}

	return rep.Err
}

func parallelLoop(l linter.Language, filteredFilepaths []string) error {
	report := make(chan linter.Report)

	toLint := l.GetLinter(CurrentLinterType)

	for _, filepath := range filteredFilepaths {
		go toLint.Run(filepath, report)
	}

	var lasterr error

	for range filteredFilepaths {
		err := printReportErr(<-report)

		if err != nil {
			// Don't just return we still have channels to join
			lasterr = err
		}
	}

	return lasterr
}
