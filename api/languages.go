package api

import "regexp"

// Directory to store default configurations for tools
var confDir = "/etc/crie" // || C:\Program Files\Common Files\crie

var lPython = language{
	name:  `python`,
	match: regexp.MustCompile(`\.py$`),
	fmt:   execCmd{`autopep8`, par{`--in-place`, `--aggressive`, `--aggressive`}, par{}},
	chk:   execCmd{`pylint`, par{`--rcfile=` + confDir + `.pylintrc`}, par{}},
}

var lPythonDeps = language{
	name:  `pythondeps`,
	match: regexp.MustCompile(`requirements.txt$`),
	fmt:   execCmd{`pur`, par{`-r`}, par{}},
}

// https://github.com/mvdan/sh/releases/download/v1.3.0/shfmt_v1.3.0_linux_amd64
// https://github.com/koalaman/shellcheck
var lBash = language{name: `bash`,
	match: regexp.MustCompile(`\.bash$`),
	fmt:   execCmd{`shfmt`, par{`-w`, `-ln`, `bash`}, par{}},
	chk:   execCmd{`shellcheck`, par{`-x`, `--shell=bash`, `-Calways`}, par{}}}

var lSh = language{name: `sh`,
	match: regexp.MustCompile(`\.sh$|/script/[^.]*$`),
	fmt:   execCmd{`shfmt`, par{`-w`, `-ln`, `posix`}, par{}},
	chk:   execCmd{`shellcheck`, par{`-x`, `--shell=sh`, `-Calways`}, par{}}}

// https://github.com/lukasmartinelli/hadolint
var lDocker = language{name: `docker`,
	match: regexp.MustCompile(`Dockerfile$`),
	chk:   execCmd{`hadolint`, par{`--ignore`, `DL3007`, `--ignore`, `DL3018`, `--ignore`, `DL3016`, `--ignore`, `DL4006`}, par{}}}

//	fmt:   execCmd{`dockfmt`, par{`fmt`, `-w`}, par{}}}

// https://github.com/adrienverge/yamllint
var lYML = language{name: `yml`,
	match: regexp.MustCompile(`\.yml$|\.yaml$`),
	chk:   execCmd{`yamllint`, par{`-c=` + confDir + `/yaml/.yamllintrc`}, par{}}}

var lTerraform = language{name: `terraform`,
	match: regexp.MustCompile(`\.tf$`),
	fmt:   execCmd{`terraform`, par{`fmt`}, par{}},
	chk:   execCmd{`terraform`, par{`fmt`, `-check=true`}, par{}}}

// https://blog.jetbrains.com/webstorm/2017/01/webstorm-2017-1-eap-171-2272/
// https://github.com/standard/standard
var lJavascript = language{name: `javascript`,
	match: regexp.MustCompile(`\.js$|\.jsx$`),
	fmt:   execCmd{`standard`, par{`--fix`}, par{}},
	chk:   execCmd{`standard`, par{}, par{}}}

// https://golang.org/cmd/gofmt/
var lGolang = language{name: `golang`,
	match: regexp.MustCompile(`\.go$`),
	fmt:   execCmd{`gofmt`, par{`-l`, `-w`}, par{}},
	chk:   execCmd{`golint`, par{`-set_exit_status`}, par{}}}

// https://github.com/wooorm/remark-lint
var lMarkdown = language{name: `markdown`,
	match: regexp.MustCompile(`\.md$`),
	fmt:   execCmd{`remark`, par{`--use`, `remark-preset-lint-recommended`}, par{`-o`}},
	chk:   execCmd{`vale`, par{`--config`, `` + confDir + `/markdown/.vale.ini`}, par{}}}

var lASCIIDoctor = language{name: `asciidoctor`,
	match: regexp.MustCompile(`\.adoc$`),
	chk:   execCmd{`vale`, par{`--config`, `` + confDir + `/markdown/.vale.ini`}, par{}}}

// https://github.com/zaach/jsonlint
var lJSON = language{name: `json`,
	match: regexp.MustCompile(`\.json$|\.JSON$`),
	fmt:   execCmd{`jsonlint`, par{`-i`, `-s`, `-c`, `-q`}, par{}},
	chk:   execCmd{`jsonlint`, par{`-q`}, par{}}}

// noExplicitConstructor and noConstructor unfortunately have problems with CUDA_CALLABLE
var lCpp = language{name: `cpp`,
	match: regexp.MustCompile(`\.cc$|\.cpp$`),
	fmt:   execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}},
	chk:   execCmd{`cppcheck`, par{`--enable=all`, `--language=c++`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}, par{}}}

var lCppheaders = language{name: `cppheaders`,
	match: regexp.MustCompile(`\.h$|\.hpp$`),
	fmt:   execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}}}

var lC = language{name: `c`,
	match: regexp.MustCompile(`\.c$`),
	fmt:   execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}},
	chk:   execCmd{`cppcheck`, par{`--enable=all`, `--language=c`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}, par{}}}

var lCmake = language{name: `cmake`,
	match: regexp.MustCompile(`CMakeLists.txt$|\.cmake$`),
	chk:   execCmd{`cmakelint`, par{`--config=` + confDir + `/cmake/.cmakelintrc`}, par{}}}

var lAnsible = language{name: `ansible`,
	match: regexp.MustCompile(`playbook.yml$`),
	chk:   execCmd{`ansible-lint`, par{}, par{}}}

var lDockerCompose = language{name: `dockercompose`,
	match: regexp.MustCompile(`docker-compose.yml$|docker-compose.yaml$`),
	chk:   execCmd{`docker-compose`, par{`-f`}, par{`config`, `-q`}}}
