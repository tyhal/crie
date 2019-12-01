package imp

import (
	"github.com/tyhal/crie/api/linter"
	"regexp"
)

// Directory to store default configurations for tools
var confDir = "/etc/crie" // || C:\Program Files\Common Files\crie

// LanguageList is a monolithic configuration of all cries standard linters
var LanguageList = []linter.Language{
	{
		Name:  `python`,
		Match: regexp.MustCompile(`\.py$`),
		Fmt:   execCmd{`autopep8`, par{`--in-place`, `--aggressive`, `--aggressive`}, par{}},
		Chk:   execCmd{`pylint`, par{`--rcfile=` + confDir + `/python/.pylintrc`}, par{}},
	}, {
		Name:  `pythondeps`,
		Match: regexp.MustCompile(`requirements.txt$`),
		Fmt:   execCmd{`pur`, par{`-r`}, par{}},
	},

	// https://github.com/mvdan/sh/releases/download/v1.3.0/shfmt_v1.3.0_linux_amd64
	// https://github.com/koalaman/shellcheck
	{
		Name:  `bash`,
		Match: regexp.MustCompile(`\.bash$`),
		Fmt:   execCmd{`shfmt`, par{`-w`, `-ln`, `bash`}, par{}},
		Chk:   execCmd{`shellcheck`, par{`-x`, `--shell=bash`, `-Calways`}, par{}}},
	{
		Name:  `sh`,
		Match: regexp.MustCompile(`\.sh$|/script/[^.]*$`),
		Fmt:   execCmd{`shfmt`, par{`-w`, `-ln`, `posix`}, par{}},
		Chk:   execCmd{`shellcheck`, par{`-x`, `--shell=sh`, `-Calways`}, par{}}},

	// https://github.com/lukasmartinelli/hadolint
	{
		Name:  `docker`,
		Match: regexp.MustCompile(`Dockerfile$`),
		Chk:   execCmd{`hadolint`, par{`--ignore`, `DL3007`, `--ignore`, `DL3018`, `--ignore`, `DL3016`, `--ignore`, `DL4006`}, par{}}},

	//	Fmt:   execCmd{`dockfmt`, par{`fmtConf`, `-w`}, par{}}}

	// https://github.com/adrienverge/yamllint
	{
		Name:  `yml`,
		Match: regexp.MustCompile(`\.yml$|\.yaml$`),
		Chk:   execCmd{`yamllint`, par{`-c=` + confDir + `/yaml/.yamllintrc`}, par{}}},

	{
		Name:  `terraform`,
		Match: regexp.MustCompile(`\.tf$`),
		Fmt:   execCmd{`terraform`, par{`fmt`}, par{}},
		Chk:   execCmd{`terraform`, par{`fmt`, `-check=true`}, par{}}},

	// https://blog.jetbrains.com/webstorm/2017/01/webstorm-2017-1-eap-171-2272/
	// https://github.com/standard/standard
	{
		Name:  `javascript`,
		Match: regexp.MustCompile(`\.js$|\.jsx$`),
		Fmt:   execCmd{`standard`, par{`--fix`}, par{}},
		Chk:   execCmd{`standard`, par{}, par{}}},

	// https://golang.org/cmd/gofmt/
	{
		Name:  `golang`,
		Match: regexp.MustCompile(`\.go$`),
		Fmt:   execCmd{`gofmt`, par{`-l`, `-w`}, par{}},
		Chk:   execCmd{`golint`, par{`-set_exit_status`}, par{}},
	},

	// https://github.com/wooorm/remark-lint
	{
		Name:  `markdown`,
		Match: regexp.MustCompile(`\.md$`),
		Fmt:   execCmd{`remark`, par{`--use`, `remark-preset-lint-recommended`}, par{`-o`}},
		Chk:   newValeLint(confDir + `/markdown/.vale.ini`)},

	{
		Name:  `asciidoctor`,
		Match: regexp.MustCompile(`\.adoc$`),
		Chk:   newValeLint(confDir + `/markdown/.vale.ini`)},

	// https://github.com/zaach/jsonlint
	{
		Name:  `json`,
		Match: regexp.MustCompile(`\.json$|\.JSON$`),
		Fmt:   execCmd{`jsonlint`, par{`-i`, `-s`, `-c`, `-q`}, par{}},
		Chk:   execCmd{`jsonlint`, par{`-q`}, par{}}},

	// noExplicitConstructor and noConstructor unfortunately have problems with CUDA_CALLABLE
	{
		Name:  `cpp`,
		Match: regexp.MustCompile(`\.cc$|\.cpp$`),
		Fmt:   execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}},
		Chk:   execCmd{`cppcheck`, par{`--enable=all`, `--language=c++`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}, par{}}},

	{
		Name:  `cppheaders`,
		Match: regexp.MustCompile(`\.h$|\.hpp$`),
		Fmt:   execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}}},

	{
		Name:  `c`,
		Match: regexp.MustCompile(`\.c$`),
		Fmt:   execCmd{`clang-format`, par{`-style=file`, `-i`}, par{}},
		Chk:   execCmd{`cppcheck`, par{`--enable=all`, `--language=c`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}, par{}}},

	{
		Name:  `cmake`,
		Match: regexp.MustCompile(`CMakeLists.txt$|\.cmake$`),
		Chk:   execCmd{`cmakelint`, par{`--config=` + confDir + `/cmake/.cmakelintrc`}, par{}}},

	{
		Name:  `ansible`,
		Match: regexp.MustCompile(`playbook.yml$`),
		Chk:   execCmd{`ansible-lint`, par{}, par{}}},

	{
		Name:  `dockercompose`,
		Match: regexp.MustCompile(`docker-compose.yml$|docker-compose.yaml$`),
		Chk:   execCmd{`docker-compose`, par{`-f`}, par{`config`, `-q`}}}}
