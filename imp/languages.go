package imp

import (
	"github.com/tyhal/crie/api"
	"github.com/tyhal/crie/api/imp"
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
		Fmt:   &imp.ExecCmd{Bin: `autopep8`, FrontPar: api.Par{`--in-place`, `--aggressive`, `--aggressive`}},
		Chk:   &imp.ExecCmd{Bin: `pylint`, FrontPar: api.Par{`--rcfile=` + confDir + `/python/.pylintrc`}},
	},
	{
		Name:  `pythondeps`,
		Match: regexp.MustCompile(`requirements.txt$`),
		Fmt:   &imp.ExecCmd{Bin: `pur`, FrontPar: api.Par{`-r`}},
	},
	{
		Name:  `proto`,
		Match: regexp.MustCompile(`.proto$`),
		Chk:   &protoLint{Fix: false},
		Fmt:   &protoLint{Fix: true},
	},

	// https://github.com/mvdan/sh/releases/download/v1.3.0/shfmt_v1.3.0_linux_amd64
	// https://github.com/koalaman/shellcheck
	{
		Name:  `bash`,
		Match: regexp.MustCompile(`\.bash$`),
		Fmt:   &imp.ExecCmd{Bin: `shfmt`, FrontPar: api.Par{`-w`, `-ln`, `bash`}},
		Chk:   &imp.ExecCmd{Bin: `shellcheck`, FrontPar: api.Par{`-x`, `--shell=bash`, `-Calways`},
			Docker: imp.DockerCmd{Image: "docker.io/tyhal/hadolint-hadolint:v1.18.0"}}},
	{
		Name:  `sh`,
		Match: regexp.MustCompile(`\.sh$|/script/[^.]*$`),
		Fmt:   &imp.ExecCmd{Bin: `shfmt`, FrontPar: api.Par{`-w`, `-ln`, `posix`}},
		Chk: &imp.ExecCmd{Bin: `shellcheck`, FrontPar: api.Par{`-x`, `--shell=sh`, `-Calways`},
			Docker: imp.DockerCmd{Image: "docker.io/tyhal/hadolint-hadolint:v1.18.0"}}},

	// https://github.com/lukasmartinelli/hadolint
	{
		Name:  `docker`,
		Match: regexp.MustCompile(`Dockerfile$`),
		Chk: &imp.ExecCmd{
			Bin:      `hadolint`,
			FrontPar: api.Par{`--ignore`, `DL3007`, `--ignore`, `DL3018`, `--ignore`, `DL3016`, `--ignore`, `DL4006`},
			Docker:   imp.DockerCmd{Image: "docker.io/tyhal/hadolint-hadolint:v1.18.0"}}},

	//	Fmt:   ExecCmd{`dockfmt`, par{`fmt`, `-w`}, par{}}}

	// https://github.com/adrienverge/yamllint
	{
		Name:  `yml`,
		Match: regexp.MustCompile(`\.yml$|\.yaml$`),
		Chk:   &imp.ExecCmd{Bin: `yamllint`, FrontPar: api.Par{`-c=` + confDir + `/yaml/.yamllintrc`}}},

	{
		Name:  `terraform`,
		Match: regexp.MustCompile(`\.tf$`),
		Fmt:   &imp.ExecCmd{Bin: `terraform`, FrontPar: api.Par{`fmt`}},
		Chk:   &imp.ExecCmd{Bin: `terraform`, FrontPar: api.Par{`fmt`, `-check=true`}}},

	// https://blog.jetbrains.com/webstorm/2017/01/webstorm-2017-1-eap-171-2272/
	// https://github.com/standard/standard
	{
		Name:  `javascript`,
		Match: regexp.MustCompile(`\.js$|\.jsx$`),
		Fmt:   &imp.ExecCmd{Bin: `standard`, FrontPar: api.Par{`--fix`}},
		Chk:   &imp.ExecCmd{Bin: `standard`}},

	// https://golang.org/cmd/gofmt/
	{
		Name:  `golang`,
		Match: regexp.MustCompile(`\.go$`),
		Fmt:   &imp.ExecCmd{Bin: `gofmt`, FrontPar: api.Par{`-l`, `-w`}},
		Chk:   &imp.ExecCmd{Bin: `golint`, FrontPar: api.Par{`-set_exit_status`}},
	},

	// https://github.com/wooorm/remark-lint
	{
		Name:  `markdown`,
		Match: regexp.MustCompile(`\.md$`),
		Fmt:   &imp.ExecCmd{Bin: `remark`, FrontPar: api.Par{`--use`, `remark-preset-lint-recommended`}, EndPar: api.Par{`-o`}},
		Chk:   newValeLint(confDir + `/markdown/.vale.ini`)},

	{
		Name:  `asciidoctor`,
		Match: regexp.MustCompile(`\.adoc$`),
		Chk:   newValeLint(confDir + `/markdown/.vale.ini`)},

	// https://github.com/zaach/jsonlint
	{
		Name:  `json`,
		Match: regexp.MustCompile(`\.json$|\.JSON$`),
		Fmt:   &imp.ExecCmd{Bin: `jsonlint`, FrontPar: api.Par{`-i`, `-s`, `-c`, `-q`}},
		Chk:   &imp.ExecCmd{Bin: `jsonlint`, FrontPar: api.Par{`-q`}}},

	// noExplicitConstructor and noConstructor unfortunately have problems with CUDA_CALLABLE
	{
		Name:  `cpp`,
		Match: regexp.MustCompile(`\.cc$|\.cpp$`),
		Fmt:   &imp.ExecCmd{Bin: `clang-format`, FrontPar: api.Par{`-style=file`, `-i`}},
		Chk:   &imp.ExecCmd{Bin: `cppcheck`, FrontPar: api.Par{`--enable=all`, `--language=c++`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}}},

	{
		Name:  `cppheaders`,
		Match: regexp.MustCompile(`\.h$|\.hpp$`),
		Fmt:   &imp.ExecCmd{Bin: `clang-format`, FrontPar: api.Par{`-style=file`, `-i`}}},

	{
		Name:  `c`,
		Match: regexp.MustCompile(`\.c$`),
		Fmt:   &imp.ExecCmd{Bin: `clang-format`, FrontPar: api.Par{`-style=file`, `-i`}},
		Chk:   &imp.ExecCmd{Bin: `cppcheck`, FrontPar: api.Par{`--enable=all`, `--language=c`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}}},

	{
		Name:  `cmake`,
		Match: regexp.MustCompile(`CMakeLists.txt$|\.cmake$`),
		Chk:   &imp.ExecCmd{Bin: `cmakelint`, FrontPar: api.Par{`--config=` + confDir + `/cmake/.cmakelintrc`}}},

	{
		Name:  `ansible`,
		Match: regexp.MustCompile(`playbook.yml$`),
		Chk:   &imp.ExecCmd{Bin: `ansible-lint`}},

	{
		Name:  `dockercompose`,
		Match: regexp.MustCompile(`docker-compose.yml$|docker-compose.yaml$`),
		Chk:   &imp.ExecCmd{Bin: `docker-compose`, FrontPar: api.Par{`-f`}, EndPar: api.Par{`config`, `-q`}}}}
