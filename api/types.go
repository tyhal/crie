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
	IsRepo          bool
	ConfPath        string
	LintType        string
	ContinueOnError bool
	Languages       []linter.Language
	GitDiff         int
	SingleLang      string
	allFiles        []string // allFiles the list of loaded files that need to be parsed
	gitFiles        []string // gitFiles the list of loaded files that 'might' need to be parsed

}

type par []string
