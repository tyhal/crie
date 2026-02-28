package runner

import (
	"regexp"
)

// Options are the core flags and settings to change execution
type Options struct {
	Continue      bool   `json:"continue,omitempty" yaml:"continue,omitempty"`
	Passes        bool   `json:"passes,omitempty" yaml:"passes,omitempty"`
	GitTarget     string `json:"gitTarget,omitempty" yaml:"gitTarget,omitempty"`
	GitDiff       bool   `json:"gitDiff,omitempty" yaml:"gitDiff,omitempty"`
	Only          string `json:"only,omitempty" yaml:"only,omitempty"`
	StrictLogging bool   `json:"-" yaml:"-"`
}

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	Options      Options
	Ignore       *regexp.Regexp
	NamedMatches NamedMatches
}
