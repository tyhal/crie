package clitool

import "regexp"

type par []string

type execCmd struct {
	bin         string
	frontparams par
	endparam    par
}

type Language struct {
	name    string
	image   string         // A docker image to use if the binary isn't on the local system
	match   *regexp.Regexp // Regex to identify files
	fmtConf execCmd        // Formatting tool
	chkConf execCmd        // Convention linting tool - Errors on any problem
}


