package conf

import (
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	"mvdan.cc/sh/v3/syntax"
	"regexp"
)

var imgHadolint = "docker.io/hadolint/hadolint:latest-alpine"
var imgTerraform = "docker.io/hashicorp/terraform:1.0.1"
var imgShellCheck = "docker.io/koalaman/shellcheck-alpine:stable"

var imgCrieNpm = "docker.io/tyhal/crie-dep-npm:latest"
var imgCriePip = "docker.io/tyhal/crie-dep-pip:latest"
var imgCrieGo = "docker.io/tyhal/crie-dep-go:latest"
var imgCrieApk = "docker.io/tyhal/crie-dep-apk:latest"

// LanguageList is a monolithic configuration of all cries standard linters
var LanguageList = []linter.Language{
	{
		Name:  `python`,
		Match: regexp.MustCompile(`\.py$`),
		Fmt:   &cli.Lint{Bin: `black`, FrontPar: cli.Par{`--quiet`}, Container: cli.Container{Image: imgCriePip}},
		Chk:   &cli.Lint{Bin: `pylint`, Container: cli.Container{Image: imgCriePip}},
	},
	{
		Name:  `proto`,
		Match: regexp.MustCompile(`.proto$`),
		Chk:   &cli.Lint{Bin: `protolint`, Container: cli.Container{Image: imgCrieGo}},
		Fmt:   &cli.Lint{Bin: `protolint`, Container: cli.Container{Image: imgCrieGo}},
	},

	// https://github.com/mvdan/sh/releases/download/v1.3.0/shfmt_v1.3.0_linux_amd64
	// https://github.com/koalaman/shellcheck
	{
		Name:  `bash`,
		Match: regexp.MustCompile(`\.bash$`),
		Fmt:   &shfmt.Lint{Language: syntax.LangBash},
		Chk:   &cli.Lint{Bin: `shellcheck`, FrontPar: cli.Par{`-x`, `--shell=bash`, `-Calways`}, Container: cli.Container{Image: imgShellCheck}}},
	{
		Name:  `sh`,
		Match: regexp.MustCompile(`\.sh$|/script/[^.]*$|^script/[^.]*$`),
		Fmt:   &shfmt.Lint{Language: syntax.LangPOSIX},
		Chk:   &cli.Lint{Bin: `shellcheck`, FrontPar: cli.Par{`-x`, `--shell=sh`, `-Calways`}, Container: cli.Container{Image: imgShellCheck}}},

	// https://github.com/lukasmartinelli/hadolint
	// TODO use config file when it can be mounted into the docker cmd too
	{
		Name:  `docker`,
		Match: regexp.MustCompile(`Dockerfile$`),
		Chk: &cli.Lint{
			Bin:       `hadolint`,
			FrontPar:  cli.Par{`--ignore`, `DL3007`, `--ignore`, `DL3018`, `--ignore`, `DL3016`, `--ignore`, `DL4006`},
			Container: cli.Container{Image: imgHadolint}}},

	//	Fmt:   Lint{`dockfmt`, par{`fmt`, `-w`}, par{}}}

	// https://github.com/adrienverge/yamllint
	{
		Name:  `yml`,
		Match: regexp.MustCompile(`\.yml$|\.yaml$`),
		Chk:   &cli.Lint{Bin: `yamllint`, Container: cli.Container{Image: imgCriePip}}},

	//{
	//	Name:  `terraform`,
	//	Match: regexp.MustCompile(`\.tf$`),
	//	Fmt:   &cli.Lint{Bin: `terraform`, Container: cli.Container{Image: imgTerraform}, FrontPar: cli.Par{`fmt`}},
	//	Chk:   &cli.Lint{Bin: `terraform`, Container: cli.Container{Image: imgTerraform}, FrontPar: cli.Par{`fmt`, `-check=true`}},
	//},

	// TODO switch to eslint
	{
		Name:  `javascript`,
		Match: regexp.MustCompile(`\.js$|\.jsx$`),
		Fmt:   &cli.Lint{Bin: `standard`, FrontPar: cli.Par{`--fix`}, Container: cli.Container{Image: imgCrieNpm}},
		Chk:   &cli.Lint{Bin: `standard`, Container: cli.Container{Image: imgCrieNpm}},
	},

	// https://golang.org/cmd/gofmt/
	{
		Name:  `golang`,
		Match: regexp.MustCompile(`\.go$`),
		Fmt:   &cli.Lint{Bin: `gofmt`, FrontPar: cli.Par{`-l`, `-w`}, Container: cli.Container{Image: imgCrieGo}},
		Chk:   &cli.Lint{Bin: `golint`, FrontPar: cli.Par{`-set_exit_status`}, Container: cli.Container{Image: imgCrieGo}},
	},

	// https://github.com/wooorm/remark-lint
	{
		Name:  `markdown`,
		Match: regexp.MustCompile(`\.md$`),
		Fmt:   &cli.Lint{Bin: `remark`, FrontPar: cli.Par{`--use`, `remark-preset-lint-recommended`}, EndPar: cli.Par{`-o`}, Container: cli.Container{Image: imgCrieNpm}},
		Chk:   &cli.Lint{Bin: `vale`, FrontPar: cli.Par{`--config='/home/standards/.config/.vale.ini`}, EndPar: cli.Par{`-o`}, Container: cli.Container{Image: imgCrieGo}},
	},

	{
		Name:  `asciidoctor`,
		Match: regexp.MustCompile(`\.adoc$`),
		Chk:   &cli.Lint{Bin: `vale`, FrontPar: cli.Par{`--config='/home/standards/.config/.vale.ini`}, EndPar: cli.Par{`-o`}, Container: cli.Container{Image: imgCrieGo}},
	},

	// https://github.com/zaach/jsonlint
	{
		Name:  `json`,
		Match: regexp.MustCompile(`\.json$|\.JSON$`),
		Fmt:   &cli.Lint{Bin: `jsonlint`, FrontPar: cli.Par{`-i`, `-s`, `-c`, `-q`}, Container: cli.Container{Image: imgCrieNpm}},
		Chk:   &cli.Lint{Bin: `jsonlint`, FrontPar: cli.Par{`-q`}, Container: cli.Container{Image: imgCrieNpm}}},

	// noExplicitConstructor and noConstructor unfortunately have problems with CUDA_CALLABLE
	{
		Name:  `cpp`,
		Match: regexp.MustCompile(`\.cc$|\.cpp$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}, Container: cli.Container{Image: imgCrieApk}},
		Chk: &cli.Lint{
			Bin: `cppcheck`,
			FrontPar: cli.Par{
				`--enable=all`, `--language=c++`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`,
			},
			Container: cli.Container{Image: imgCrieApk},
		}},

	{
		Name:  `cppheaders`,
		Match: regexp.MustCompile(`\.h$|\.hpp$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}, Container: cli.Container{Image: imgCrieApk}}},

	{
		Name:  `c`,
		Match: regexp.MustCompile(`\.c$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}, Container: cli.Container{Image: imgCrieApk}},
		Chk: &cli.Lint{
			Bin: `cppcheck`,
			FrontPar: cli.Par{
				`--enable=all`, `--language=c`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`,
			},
			Container: cli.Container{Image: imgCrieApk},
		}},

	{
		Name:  `cmake`,
		Match: regexp.MustCompile(`CMakeLists.txt$|\.cmake$`),
		Chk:   &cli.Lint{Bin: `cmakelint`, FrontPar: cli.Par{"--config=/home/standards/.config/cmakelintrc"}, Container: cli.Container{Image: imgCriePip}}},

	// TODO Review tools that parse child files - ansiblelint needs to install dependencies similiar to how clang-tidy does
	//{
	//	Name:  `ansible`,
	//	Match: regexp.MustCompile(`playbook.yml$`),
	//	Chk:   &cli.Lint{Bin: `ansible-lint`, Container: cli.Container{Image: imgCriePip}},
	//},

	// TODO use v2 with go
	//{
	//	Name:  `dockercompose`,
	//	Match: regexp.MustCompile(`docker-compose.yml$|docker-compose.yaml$`),
	//	Chk:   &cli.Lint{Bin: `docker-compose`, FrontPar: cli.Par{`-f`}, EndPar: cli.Par{`config`, `-q`}}}
}
