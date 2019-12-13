package imp

// github.com/errata-ai/vale@v1.7.1

import (
	"bytes"
	"github.com/tyhal/crie/api/linter"
	"github.com/yoheimuta/protolint/api"
	"math"
)

type protoLint struct {
	Fix bool
}

func (e protoLint) Name() string {
	return "protolint"
}

func (e protoLint) WillRun() (err error) {
	return
}

func (e protoLint) MaxConcurrency() int {
	return math.MaxInt32
}

func (e protoLint) Run(filepath string, rep chan linter.Report) {
	var outB, errB bytes.Buffer
	err := api.Lint(
		filepath,
		e.Fix,
		&outB,
		&errB,
	)
	rep <- linter.Report{filepath, err, &outB, &errB}
}

