package dockfmt

import (
	"context"

	"github.com/tyhal/crie/pkg/linter"
)

const defaultIndent = 4

// LintDockFmt is a linter that uses the reteps/dockerfmt library to format dockerfiles
type LintDockFmt struct {
	Type            string `json:"type" yaml:"type" jsonschema:"a built in Docker formatter thanks to reteps (Pete Stenger)"`
	IndentSize      uint   `json:"indent_size" yaml:"indent_size" jsonschema:"Number of spaces to use for indentation"`
	TrailingNewline bool   `json:"trailing_newline" yaml:"trailing_newline" jsonschema:"End the file with a trailing newline"`
	SpaceRedirects  bool   `json:"space_redirects" yaml:"space_redirects" jsonschema:"Redirect operators will be followed by a space"`
}

var _ linter.Linter = (*LintDockFmt)(nil)

// Name returns the name of the linter
func (l *LintDockFmt) Name() string {
	return "dockfmt"
}

// Setup returns nil as there are no external deps
func (l *LintDockFmt) Setup(_ context.Context) (err error) {
	return nil
}

// Cleanup removes any additional resources created in the process
func (l *LintDockFmt) Cleanup(_ context.Context) error {
	return nil
}

// Run will just return the configured error as a report
func (l *LintDockFmt) Run(filepath string) linter.Report {
	err := l.format(filepath)
	return linter.Report{Target: filepath, Err: err, StdOut: nil, StdErr: nil}
}
