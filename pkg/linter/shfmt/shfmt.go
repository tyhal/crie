package shfmt

import (
	"bytes"
	"fmt"
	"github.com/tyhal/crie/pkg/crie/linter"
	"math"
	"mvdan.cc/sh/v3/syntax"
	"strings"
	"sync"
)

const ()

// LintShfmt contains all the information needed to configure shfmt
type LintShfmt struct {
	Type     string `json:"type" yaml:"type" jsonschema:"enum=shfmt" jsonschema_description:"a built in shell formatter thanks to mvdan"`
	Language string `json:"language" yaml:"language" jsonschema:"enum=bash,enum=posix,enum=sh,enum=mksh"`
}

// Name returns the name of the linter
func (l *LintShfmt) Name() string {
	return "shfmt"
}

// WillRun do nothing as there are no external deps
func (l *LintShfmt) WillRun() (err error) {
	return nil
}

// Cleanup removes any additional resources created in the process
func (l *LintShfmt) Cleanup(wg *sync.WaitGroup) {
	wg.Done()
}

// MaxConcurrency return max number of parallel files to fmt
func (l *LintShfmt) MaxConcurrency() int {
	return math.MaxInt32
}

// Run shfmt -w
func (l *LintShfmt) Run(filepath string, rep chan linter.Report) {
	var outB, errB bytes.Buffer

	currFmt := shfmt{
		syntax.NewParser(syntax.KeepComments(true)),
		syntax.NewPrinter(),
	}

	var lang syntax.LangVariant

	switch strings.ToLower(l.Language) {
	case "bash":
		lang = syntax.LangBash
	case "posix", "sh":
		lang = syntax.LangPOSIX
	case "mksh":
		lang = syntax.LangMirBSDKorn
	default:
		err := fmt.Errorf("unknown language variant: %s", l.Language)
		rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: &errB}
		return
	}

	syntax.Variant(lang)(currFmt.parser)
	syntax.Indent(0)(currFmt.printer)
	syntax.BinaryNextLine(false)(currFmt.printer)
	syntax.SwitchCaseIndent(false)(currFmt.printer)
	syntax.SpaceRedirects(false)(currFmt.printer)
	syntax.KeepPadding(false)(currFmt.printer)
	syntax.FunctionNextLine(false)(currFmt.printer)

	err := currFmt.formatPath(filepath, true)
	rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: &errB}
}
