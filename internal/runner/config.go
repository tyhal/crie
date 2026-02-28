package runner

import (
	"regexp"
)

// Options are the core flags and settings to change execution
type Options struct {
	Continue      bool   `json:"continue,omitempty" yaml:"continue"`
	Passes        bool   `json:"passes,omitempty" yaml:"passes"`
	GitTarget     string `json:"gitTarget,omitempty" yaml:"gitTarget"`
	GitDiff       bool   `json:"gitDiff,omitempty" yaml:"gitDiff"`
	Only          string `json:"only,omitempty" yaml:"only"`
	StrictLogging bool   `json:"-" yaml:"-"`
}

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	Options      Options
	Ignore       *regexp.Regexp
	NamedMatches NamedMatches
}
