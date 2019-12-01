package linter

import (
	"log"
	"regexp"
)

// Report is used to state what issues a given file has
type Report struct {
	File   string
	Err    error
	StdOut string
	StdErr string
}

// Linter is a simple inteface to enable a setup and check using WillRun before executing multiple Run's
type Linter interface {
	WillRun() error
	Run(filepath string, rep chan Report)
}

// Language is used to associate a file pattern to the relevant tools to check and format
type Language struct {
	Name  string
	Match *regexp.Regexp // Regex to identify files
	Fmt   Linter         // Formatting tool
	Chk   Linter         // Convention linting tool - Errors on any problem
}

// GetLinter returns the function to run when executing
// TODO remove requirement for this function
func (l Language) GetLinter(which string) Linter {
	if which == "fmt" {
		return l.Fmt
	} else if which == "chk" {
		return l.Chk
	}
	// XXX should really pass back down
	log.Fatal("No linter found '" + which + "'")
	return nil
}
