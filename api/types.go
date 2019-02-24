package api

import (
	"regexp"
)

type par []string

type execCmd struct {
	bin         string
	frontparams par
	endparam    par
}

type language struct {
	name  string
	image string         // A docker image to use if the binary isn't on the local system
	match *regexp.Regexp // Regex to identify files
	tacks bool           // If the file should use '-' in its name
	fmt   execCmd        // Formatting tool
	chk   execCmd        // Convention linting tool - Errors on any problem
}

type report struct {
	file   string
	err    error
	stdout string
	stderr string
}

type conf struct {
	Ignore   []string `yaml:"ignore"`
	ProjDirs []string `yaml:"proj_dirs"`
}

type state struct {
	IsRepo   bool
	ConfName string
}

type execselec func(language) execCmd
