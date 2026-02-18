package runner

import (
	"regexp"
)

// Options are the core flags and settings to change execution
type Options struct {
	Continue      bool   `json:"continue" yaml:"continue"`
	Passes        bool   `json:"passes" yaml:"passes"`
	GitTarget     string `json:"gitTarget" yaml:"gitTarget"`
	GitDiff       bool   `json:"gitDiff" yaml:"gitDiff"`
	Only          string `json:"only" yaml:"only"`
	StrictLogging bool   `json:"-" yaml:"-"`
}

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	Options      Options
	Ignore       *regexp.Regexp
	NamedMatches NamedMatches
}
