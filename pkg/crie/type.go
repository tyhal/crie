package crie

import (
	"github.com/tyhal/crie/pkg/crie/linter"
	"regexp"
)

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	lintType        string
	ContinueOnError bool
	StrictLogging   bool
	ShowPasses      bool
	Languages       map[string]*linter.Language
	Ignore          *regexp.Regexp
	GitDiff         int
	SingleLang      string
	fileList        []string
}
