package shfmt

import (
	"bytes"
	"fmt"
	"github.com/tyhal/crie/pkg/crie/linter"
	"gopkg.in/yaml.v3"
	"math"
	"mvdan.cc/sh/v3/syntax"
	"strings"
)

// Lint contains all the information needed to configure shfmt
type Lint struct {
	Language syntax.LangVariant `yaml:"language"`
}

// UnmarshalYAML implements custom YAML unmarshalling
func (l *Lint) UnmarshalYAML(value *yaml.Node) error {
	var temp struct {
		Language string `yaml:"language"`
	}

	if err := value.Decode(&temp); err != nil {
		return err
	}

	switch strings.ToLower(temp.Language) {
	case "bash":
		l.Language = syntax.LangBash
	case "posix", "sh":
		l.Language = syntax.LangPOSIX
	case "mksh":
		l.Language = syntax.LangMirBSDKorn
	default:
		return fmt.Errorf("unknown language variant: %s", temp.Language)
	}

	return nil
}

// Name returns the name of the linter
func (l *Lint) Name() string {
	return "shfmt"
}

// WillRun do nothing as there are no external deps
func (l *Lint) WillRun() (err error) {
	return nil
}

// Cleanup removes any additional resources created in the process
func (l *Lint) Cleanup() {
	return
}

// WaitForCleanup Useful for when Cleanup is running in the background
func (l *Lint) WaitForCleanup() (err error) {
	return nil
}

// MaxConcurrency return max number of parallel files to fmt
func (l *Lint) MaxConcurrency() int {
	return math.MaxInt32
}

// Run shfmt -w
func (l *Lint) Run(filepath string, rep chan linter.Report) {
	var outB, errB bytes.Buffer

	currFmt := shfmt{
		syntax.NewParser(syntax.KeepComments(true)),
		syntax.NewPrinter(),
	}

	syntax.Variant(l.Language)(currFmt.parser)
	syntax.Indent(0)(currFmt.printer)
	syntax.BinaryNextLine(false)(currFmt.printer)
	syntax.SwitchCaseIndent(false)(currFmt.printer)
	syntax.SpaceRedirects(false)(currFmt.printer)
	syntax.KeepPadding(false)(currFmt.printer)
	syntax.FunctionNextLine(false)(currFmt.printer)

	err := currFmt.formatPath(filepath, true)

	rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: &errB}
}
