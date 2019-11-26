package api

// TODO(tyler) this has repeated exec calls

import (
	"fmt"
	"github.com/tyhal/crie/api/linter"
	"log"
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
		if Verbose && !Quiet {
			fmt.Println("	" + rep.StdOut)
		}
		return nil
	}

	fmt.Println(" ✖️  " + rep.File)

	if Quiet {
		return rep.Err
	}

	fmt.Println("	std out : ")
	fmt.Println(rep.StdOut)
	fmt.Println("	std err : ")
	fmt.Println(rep.StdErr)

	if ContinueOnError {
		fmt.Println("	cmd err : ")
		fmt.Println("	" + rep.Err.Error())
	} else {
		log.Fatal(rep.Err)
	}

	return rep.Err
}

func parallelLoop(l linter.Linter,filteredFilepaths []string) error {
	report := make(chan linter.Report)

	for _, filepath := range filteredFilepaths {
		go l.Chk(filepath, report)
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
