package proto

import (
	"bytes"
	linter "github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/protolint/api"
	"math"
)

// Lint Lint
type Lint struct {
	Fix bool
}

// Name Name
func (e *Lint) Name() string {
	return "protolint"
}

// WillRun WillRun
func (e *Lint) WillRun() (err error) {
	return
}

// DidRun DidRun
func (e *Lint) DidRun() {
	return
}

// MaxConcurrency MaxConcurrency
func (e *Lint) MaxConcurrency() int {
	return math.MaxInt32
}

// Run run
func (e *Lint) Run(filepath string, rep chan linter.Report) {
	var outB, errB bytes.Buffer
	err := api.Lint(
		filepath,
		e.Fix,
		&outB,
		&errB,
	)
	rep <- linter.Report{File: filepath, Err: err, StdOut: &outB, StdErr: &errB}
}
