package vale

// github.com/errata-ai/vale@v1.7.1

import (
	"github.com/errata-ai/vale/check"
	"github.com/errata-ai/vale/core"
	"github.com/errata-ai/vale/lint"
	"github.com/tyhal/crie/pkg/crie/linter"
	"io"
	"math"
)

// Lint Lint
type Lint struct {
	configPath string
	linter     *lint.Linter
}

// Name Name
func (e *Lint) Name() string {
	return "vale"
}

// WillRun WillRun
func (e *Lint) WillRun() (err error) {
	config := core.NewConfig()
	config, err = core.LoadConfig(config, e.configPath, "warning", false)
	e.linter.Config = config
	e.linter.CheckManager = check.NewManager(config)
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

// Run Run
func (e *Lint) Run(filepath string, rep chan linter.Report) {
	var stdout io.Reader
	linted, err := e.linter.LintString(filepath)
	if err == nil {
		stdout, err = GetVerboseAlerts(linted, e.linter.Config.Wrap)
	}
	rep <- linter.Report{File: filepath, Err: err, StdOut: stdout}
}

// NewValeLint NewValeLint
func NewValeLint(confpath string) *Lint {
	return &Lint{confpath, &lint.Linter{}}
}
