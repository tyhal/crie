package dockfmt

import (
	"math"
	"os"
	"strings"
	"sync"

	"github.com/reteps/dockerfmt/lib"
	"github.com/tyhal/crie/internal/errchain"
	"github.com/tyhal/crie/pkg/linter"
)

const defaultIndent = 4

// LintDockFmt
type LintDockFmt struct {
	Type            string `json:"type" yaml:"type" jsonschema:"enum=dockfmt" jsonschema_description:"a built in Docker formatter thanks to reteps (Pete Stenger)"`
	IndentSize      uint   `json:"indent_size" yaml:"indent_size" jsonschema_description:"Number of spaces to use for indentation"`
	TrailingNewline bool   `json:"trailing_newline" yaml:"trailing_newline" jsonschema_description:"End the file with a trailing newline"`
	SpaceRedirects  bool   `json:"space_redirects" yaml:"space_redirects" jsonschema_description:"Redirect operators will be followed by a space"`
}

var _ linter.Linter = (*LintDockFmt)(nil)

// Name returns the name of the linter
func (l *LintDockFmt) Name() string {
	return "dockfmt"
}

// WillRun just returns the configured error
func (l *LintDockFmt) WillRun() (err error) {
	return nil
}

// Cleanup removes any additional resources created in the process
func (l *LintDockFmt) Cleanup(group *sync.WaitGroup) {
	group.Done()
}

// MaxConcurrency return max number of parallel files to fmt
func (l *LintDockFmt) MaxConcurrency() int {
	return math.MaxInt32
}

func (l *LintDockFmt) format(filepath string) error {
	if l.IndentSize == 0 {
		l.IndentSize = defaultIndent
	}
	config := &lib.Config{
		IndentSize:      l.IndentSize,
		TrailingNewline: l.TrailingNewline,
		SpaceRedirects:  l.SpaceRedirects,
	}

	src, err := os.ReadFile(filepath)
	if err != nil {
		return errchain.From(err).LinkF("Failed to read from file %s", filepath)
	}

	lines := strings.SplitAfter(strings.TrimSuffix(string(src), "\n"), "\n")
	dst := lib.FormatFileLines(lines, config)

	err = os.WriteFile(filepath, []byte(dst), 0644)
	if err != nil {
		return errchain.From(err).LinkF("Failed to write to file %s", filepath)
	}

	return nil
}

// Run will just return the configured error as a report
func (l *LintDockFmt) Run(filepath string, rep chan linter.Report) {
	err := l.format(filepath)
	rep <- linter.Report{File: filepath, Err: err, StdOut: nil, StdErr: nil}
}
