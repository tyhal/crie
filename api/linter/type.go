package linter

import (
	"log"
	"regexp"
)

type Report struct {
	File   string
	Err    error
	StdOut string
	StdErr string
}

type Linter interface {
	WillRun() error
	Run(filepath string, rep chan Report)
}

type Language struct {
	Name  string
	Match *regexp.Regexp // Regex to identify files
	Fmt   Linter         // Formatting tool
	Chk   Linter         // Convention linting tool - Errors on any problem
}

type Lint func(file string, rep chan Report)

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
