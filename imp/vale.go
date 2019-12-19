package imp

// github.com/errata-ai/vale@v1.7.1

import (
	"github.com/errata-ai/vale/check"
	"github.com/errata-ai/vale/core"
	"github.com/errata-ai/vale/lint"
	"github.com/tyhal/crie/api/linter"
	"github.com/tyhal/crie/imp/printer"
	"io"
	"math"
)

type valeLint struct {
	configPath string
	linter     *lint.Linter
}

func (e *valeLint) Name() string {
	return "vale"
}

func (e *valeLint) WillRun() (err error) {
	config := core.NewConfig()
	config, err = core.LoadConfig(config, e.configPath, "warning", false)
	e.linter.Config = config
	e.linter.CheckManager = check.NewManager(config)
	return
}

func (e *valeLint) DidRun() {
	return
}

func (e *valeLint) MaxConcurrency() int {
	return math.MaxInt32
}

func (e *valeLint) Run(filepath string, rep chan linter.Report) {
	var stdout io.Reader
	linted, err := e.linter.LintString(filepath)
	if err == nil {
		stdout, err = printer.GetVerboseAlerts(linted, e.linter.Config.Wrap)
	}
	rep <- linter.Report{File: filepath, Err: err, StdOut: stdout}
}

func newValeLint(confpath string) *valeLint {
	return &valeLint{confpath, &lint.Linter{}}
}
