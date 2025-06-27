package linter

import (
	log "github.com/sirupsen/logrus"
	"io"
	"regexp"
)

// Report is used to state what issues a given file has
type Report struct {
	File   string
	Err    error
	StdOut io.Reader
	StdErr io.Reader
}

// Linter is a simple interface to enable a setup and check using WillRun before executing multiple Run's
type Linter interface {
	Name() string
	WillRun() error
	Cleanup()
	MaxConcurrency() int
	Run(filePath string, rep chan Report)
	WaitForCleanup() error
}

// Language is used to associate a file pattern with the relevant tools to check and format
type Language struct {
	Name  string   `yaml:"name"`
	Match []string `yaml:"match"`
	Regex *regexp.Regexp
	Fmt   Linter `yaml:"fmt,omitempty"`   // Formatting tool
	Chk   Linter `yaml:"check,omitempty"` // Convention linting tool - Errors on any problem
}

// GetLinter allows for string indexing to get fmt or chk...
// TODO remove requirement for this function
func (l *Language) GetLinter(which string) Linter {
	if which == "fmt" {
		return l.Fmt
	} else if which == "chk" {
		return l.Chk
	}
	// XXX should really pass back down
	log.Fatal("No linter found '" + which + "'")
	return nil
}
