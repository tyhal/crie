package runner

import (
	"regexp"
)

// Options are the core flags and settings to change execution
type Options struct {
	Continue      bool   `json:"continue,omitzero" yaml:"continue"`
	Passes        bool   `json:"passes,omitzero" yaml:"passes"`
	GitTarget     string `json:"gitTarget,omitzero" yaml:"gitTarget"`
	GitDiff       bool   `json:"gitDiff,omitzero" yaml:"gitDiff"`
	Only          string `json:"only,omitzero" yaml:"only"`
	StrictLogging bool   `json:"-" yaml:"-"`
}

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	Options      Options
	Ignore       *regexp.Regexp
	NamedMatches NamedMatches
}
