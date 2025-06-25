package shfmt

import (
	"bytes"
	"github.com/tyhal/crie/pkg/crie/linter"
	"math"
	"mvdan.cc/sh/v3/syntax"
)

// Lint contains all the information needed to configure shfmt
type Lint struct {
	Language syntax.LangVariant
}

// Name return the name of the linter
func (e *Lint) Name() string {
	return "shfmt"
}

// WillRun do nothing as there are no external deps
func (e *Lint) WillRun() (err error) {
	return nil
}

// Cleanup remove any additional resources created in the process
func (e *Lint) Cleanup() {
	return
}

// WaitForCleanup Useful for when Cleanup is running in the background
func (e *Lint) WaitForCleanup() (err error) {
	return nil
}

// MaxConcurrency return max amount of parallel files to fmt
func (e *Lint) MaxConcurrency() int {
	return math.MaxInt32
}

// Run run shfmt -w
func (e *Lint) Run(filepath string, rep chan linter.Report) {
	var outB, errB bytes.Buffer

	currFmt := shfmt{
		syntax.NewParser(syntax.KeepComments(true)),
		syntax.NewPrinter(),
	}

	syntax.Variant(e.Language)(currFmt.parser)
	syntax.Indent(0)(currFmt.printer)
	syntax.BinaryNextLine(false)(currFmt.printer)
	syntax.SwitchCaseIndent(false)(currFmt.printer)
	syntax.SpaceRedirects(false)(currFmt.printer)
	syntax.KeepPadding(false)(currFmt.printer)
	syntax.FunctionNextLine(false)(currFmt.printer)

	err := currFmt.formatPath(filepath, true)

	rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: &errB}
}
