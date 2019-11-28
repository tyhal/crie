package imp

import (
	"errors"
	"github.com/tyhal/crie/api/linter"
)

type valeLint struct {
	config string
}

func (e valeLint) WillRun() error {
	return errors.New("not implemented")
}

func (e valeLint) Run(filepath string, rep chan linter.Report) {
	err := errors.New("not implemented")

	rep <- linter.Report{filepath, err, "stdout", "stderr"}
}
