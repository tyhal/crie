package crie

import (
	"github.com/tyhal/crie/pkg/crie/linter"
)

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	ConfPath        string
	lintType        string
	ContinueOnError bool
	StrictLogging   bool
	ShowPasses      bool
	Languages       []linter.Language
	GitDiff         int
	SingleLang      string
	fileList        []string
}

// ProjectSettings is the configuration for users to override defaults and ignore files
type ProjectSettings struct {
	Ignore   []string `yaml:"ignore"`
	ProjDirs []string `yaml:"proj_dirs"`
}
