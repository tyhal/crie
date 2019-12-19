package api

import (
	"github.com/tyhal/crie/api/linter"
)

// fileSettings simply adjusts what we include in a normal lint
type fileSettings struct {
	Ignore   []string `yaml:"ignore"`
	ProjDirs []string `yaml:"proj_dirs"`
}

// ProjectLintConfiguration is what is required for an entire project to be linted
type ProjectLintConfiguration struct {
	ConfPath        string
	lintType        string
	ContinueOnError bool
	ShowPasses      bool
	Languages       []linter.Language
	GitDiff         int
	SingleLang      string
	fileList        []string
}

// Par represents cli parameters
type Par []string
