package native

import "regexp"

var lPythonDeps = language{
	name:    `pythondeps`,
	match:   regexp.MustCompile(`requirements.txt$`),
	fmtConf: execCmd{`pur`, par{`-r`}, par{}},
}

