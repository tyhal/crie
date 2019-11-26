package clitool

import (
	"regexp"
)

// Directory to store default configurations for tools
var confDir = "/etc/crie" // || C:\Program Files\Common Files\crie

var LanguageList = []Language{
	{
		name:    `python`,
		match:   regexp.MustCompile(`\.py$`),
		fmtConf: execCmd{`autopep8`, par{`--in-place`, `--aggressive`, `--aggressive`}, par{}},
		chkConf: execCmd{`pylint`, par{`--rcfile=` + confDir + `/python/.pylintrc`}, par{}},
	}, {
		name:    `pythondeps`,
		match:   regexp.MustCompile(`requirements.txt$`),
		fmtConf: execCmd{`pur`, par{`-r`}, par{}},
	},

	// https://github.com/mvdan/sh/releases/download/v1.3.0/shfmt_v1.3.0_linux_amd64
	// https://github.com/koalaman/shellcheck
	{
		name:    `bash`,
		match:   regexp.MustCompile(`\.bash$`),
		fmtConf: execCmd{`shfmt`, par{`-w`, `-ln`, `bash`}, par{}},
		chkConf: execCmd{`shellcheck`, par{`-x`, `--shell=bash`, `-Calways`}, par{}}},
	{name: `sh`,
		match:   regexp.MustCompile(`\.sh$|/script/[^.]*$`),
		fmtConf: execCmd{`shfmt`, par{`-w`, `-ln`, `posix`}, par{}},
		chkConf: execCmd{`shellcheck`, par{`-x`, `--shell=sh`, `-Calways`}, par{}}},

	// https://github.com/lukasmartinelli/hadolint
	{name: `docker`,
		match:   regexp.MustCompile(`Dockerfile$`),
		chkConf: execCmd{`hadolint`, par{`--ignore`, `DL3007`, `--ignore`, `DL3018`, `--ignore`, `DL3016`, `--ignore`, `DL4006`}, par{}}},

	//	fmtConf:   execCmd{`dockfmt`, par{`fmtConf`, `-w`}, par{}}}

	// https://github.com/adrienverge/yamllint
	{name: `yml`,
		match:   regexp.MustCompile(`\.yml$|\.yaml$`),
		chkConf: execCmd{`yamllint`, par{`-c=` + confDir + `/yaml/.yamllintrc`}, par{}}},

	{name: `terraform`,
		match:   regexp.MustCompile(`\.tf$`),
		fmtConf: execCmd{`terraform`, par{`fmtConf`}, par{}},
		chkConf: execCmd{`terraform`, par{`fmtConf`, `-check=true`}, par{}}},

	// https://blog.jetbrains.com/webstorm/2017/01/webstorm-2017-1-eap-171-2272/
	// https://github.com/standard/standard
	{name: `javascript`,
		match:   regexp.MustCompile(`\.js$|\.jsx$`),
		fmtConf: execCmd{`standard`, par{`--fix`}, par{}},
		chkConf: execCmd{`standard`, par{}, par{}}},

	// https://golang.org/cmd/gofmt/
	{name: `golang`,
		match:   regexp.MustCompile(`\.go$`),
		fmtConf: execCmd{`gofmt`, par{`-l`, `-w`}, par{}},
		chkConf: execCmd{`golint`, par{`-set_exit_status`}, par{}}},

	// https://github.com/wooorm/remark-lint
	{name: `markdown`,
		match:   regexp.MustCompile(`\.md$`),
		fmtConf: execCmd{`remark`, par{`--use`, `remark-preset-lint-recommended`}, par{`-o`}},
		chkConf: execCmd{`vale`, par{`--config`, `` + confDir + `/markdown/.vale.ini`}, par{}}},

	{name: `asciidoctor`,
		match:   regexp.MustCompile(`\.adoc$`),
		chkConf: execCmd{`vale`, par{`--config`, `` + confDir + `/markdown/.vale.ini`}, par{}}},

	// https://github.com/zaach/jsonlint
	{name: `json`,
		match:   regexp.MustCompile(`\.json$|\.JSON$`),
		fmtConf: execCmd{`jsonlint`, par{`-i`, `-s`, `-c`, `-q`}, par{}},
		chkConf: execCmd{`jsonlint`, par{`-q`}, par{}}},

	// noExplicitConstructor and noConstructor unfortunately have problems with CUDA_CALLABLE
	{name: `cpp`,
		match:   regexp.MustCompile(`\.cc$|\.cpp$`),
		fmtConf: execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}},
		chkConf: execCmd{`cppcheck`, par{`--enable=all`, `--Language=c++`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}, par{}}},

	{name: `cppheaders`,
		match:   regexp.MustCompile(`\.h$|\.hpp$`),
		fmtConf: execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}}},

	{name: `c`,
		match:   regexp.MustCompile(`\.c$`),
		fmtConf: execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}},
		chkConf: execCmd{`cppcheck`, par{`--enable=all`, `--Language=c`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}, par{}}},

	{name: `cmake`,
		match:   regexp.MustCompile(`CMakeLists.txt$|\.cmake$`),
		chkConf: execCmd{`cmakelint`, par{`--config=` + confDir + `/cmake/.cmakelintrc`}, par{}}},

	{name: `ansible`,
		match:   regexp.MustCompile(`playbook.yml$`),
		chkConf: execCmd{`ansible-lint`, par{}, par{}}},

	{name: `dockercompose`,
		match:   regexp.MustCompile(`docker-compose.yml$|docker-compose.yaml$`),
		chkConf: execCmd{`docker-compose`, par{`-f`}, par{`config`, `-q`}}}}
