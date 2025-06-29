package linter

import (
	"fmt"
	"io"
	"regexp"
	"sync"
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
	Cleanup(wg *sync.WaitGroup)
	MaxConcurrency() int
	Run(filePath string, rep chan Report)
}

// Language is used to associate a file pattern with the relevant tools to check and format
type Language struct {
	Regex *regexp.Regexp
	Fmt   Linter
	Chk   Linter
}

// GetLinter allows for string indexing to get fmt or chk...
// TODO remove requirement for this function
func (l *Language) GetLinter(which string) (Linter, error) {

	switch which {
	case "fmt":
		return l.Fmt, nil
	case "chk":
		return l.Chk, nil
	}

	return nil, fmt.Errorf("no linter found %s", which)
}
