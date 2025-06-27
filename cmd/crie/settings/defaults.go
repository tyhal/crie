package settings

import (
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	"mvdan.cc/sh/v3/syntax"
	"regexp"
)

// 	_ "embed"
// //go:embed defaults.yml
// var languagesYAML []byte

// TODO Move out: this isn't actually part of the CLI and is an example

var imgHadolint = "docker.io/hadolint/hadolint:latest-alpine"
var imgTerraform = "docker.io/hashicorp/terraform:1.3.5"
var imgTflint = "ghcr.io/terraform-linters/tflint"
var imgShellCheck = "docker.io/koalaman/shellcheck-alpine:stable"

var imgCrieNpm = "docker.io/tyhal/crie-dep-npm:latest"
var imgCriePip = "docker.io/tyhal/crie-dep-pip:latest"
var imgCrieGo = "docker.io/tyhal/crie-dep-go:latest"
var imgCrieApk = "docker.io/tyhal/crie-dep-apk:latest"

// Configuration is a monolithic configuration of all cries standard linters
var Configuration = []linter.Language{
	{
		Name:  `python`,
		Regex: regexp.MustCompile(`\.py$`),
		Fmt:   &cli.Lint{Bin: `black`, FrontPar: cli.Par{`--quiet`}, ContainerImage: imgCriePip},
		Chk:   &cli.Lint{Bin: `pylint`, ContainerImage: imgCriePip},
	},
	{
		Name:  `proto`,
		Regex: regexp.MustCompile(`.proto$`),
		Chk:   &cli.Lint{Bin: `protolint`, ContainerImage: imgCrieGo},
		Fmt:   &cli.Lint{Bin: `protolint`, ContainerImage: imgCrieGo},
	},

	// https://github.com/mvdan/sh/releases/download/v1.3.0/shfmt_v1.3.0_linux_amd64
	// https://github.com/koalaman/shellcheck
	{
		Name:  `bash`,
		Regex: regexp.MustCompile(`\.bash$`),
		Fmt:   &shfmt.Lint{Language: syntax.LangBash},
		Chk:   &cli.Lint{Bin: `shellcheck`, FrontPar: cli.Par{`-x`, `--shell=bash`, `-Calways`}, ContainerImage: imgShellCheck}},
	{
		Name:  `sh`,
		Regex: regexp.MustCompile(`\.sh$|/script/[^.]*$|^script/[^.]*$`),
		Fmt:   &shfmt.Lint{Language: syntax.LangPOSIX},
		Chk:   &cli.Lint{Bin: `shellcheck`, FrontPar: cli.Par{`-x`, `--shell=sh`, `-Calways`}, ContainerImage: imgShellCheck}},

	// https://github.com/lukasmartinelli/hadolint
	// TODO use config file when it can be mounted into the docker cmd too
	{
		Name:  `docker`,
		Regex: regexp.MustCompile(`Dockerfile$`),
		Chk: &cli.Lint{
			Bin:            `hadolint`,
			FrontPar:       cli.Par{`--ignore`, `DL3007`, `--ignore`, `DL3018`, `--ignore`, `DL3016`, `--ignore`, `DL4006`},
			ContainerImage: imgHadolint}},

	//	Fmt:   Lint{`dockfmt`, par{`fmt`, `-w`}, par{}}}

	// https://github.com/adrienverge/yamllint
	{
		Name:  `yml`,
		Regex: regexp.MustCompile(`\.yml$|\.yaml$`),
		Chk:   &cli.Lint{Bin: `yamllint`, ContainerImage: imgCriePip}},

	{
		Name:  `terraform`,
		Regex: regexp.MustCompile(`\.tf$`),
		Fmt:   &cli.Lint{Bin: `terraform`, FrontPar: cli.Par{`fmt`}, ContainerImage: imgTerraform},
		Chk:   &cli.Lint{Bin: `tflint`, FrontPar: cli.Par{`--filter`}, ChDir: true, ContainerImage: imgTflint},
	},

	// TODO switch to eslint
	{
		Name:  `javascript`,
		Regex: regexp.MustCompile(`\.js$|\.jsx$`),
		Fmt:   &cli.Lint{Bin: `standard`, FrontPar: cli.Par{`--fix`}, ContainerImage: imgCrieNpm},
		Chk:   &cli.Lint{Bin: `standard`, ContainerImage: imgCrieNpm},
	},

	// https://golang.org/cmd/gofmt/
	{
		Name:  `golang`,
		Regex: regexp.MustCompile(`\.go$`),
		Fmt:   &cli.Lint{Bin: `gofmt`, FrontPar: cli.Par{`-l`, `-w`}, ContainerImage: imgCrieGo},
		Chk:   &cli.Lint{Bin: `golint`, FrontPar: cli.Par{`-set_exit_status`}, ContainerImage: imgCrieGo},
	},

	// https://github.com/wooorm/remark-lint
	{
		Name:  `markdown`,
		Regex: regexp.MustCompile(`\.md$`),
		Fmt:   &cli.Lint{Bin: `remark`, FrontPar: cli.Par{`--use`, `remark-preset-lint-recommended`}, EndPar: cli.Par{`-o`}, ContainerImage: imgCrieNpm},
		Chk:   &cli.Lint{Bin: `vale`, FrontPar: cli.Par{`--config=/etc/vale/.vale.ini`}, EndPar: cli.Par{}, ContainerImage: imgCrieGo},
	},

	{
		Name:  `asciidoctor`,
		Regex: regexp.MustCompile(`\.adoc$`),
		Chk:   &cli.Lint{Bin: `vale`, FrontPar: cli.Par{`--config=/etc/vale/.vale.ini`}, EndPar: cli.Par{}, ContainerImage: imgCrieGo},
	},

	// https://github.com/zaach/jsonlint
	{
		Name:  `json`,
		Regex: regexp.MustCompile(`\.json$|\.JSON$`),
		Fmt:   &cli.Lint{Bin: `jsonlint`, FrontPar: cli.Par{`-i`, `-s`, `-c`, `-q`}, ContainerImage: imgCrieNpm},
		Chk:   &cli.Lint{Bin: `jsonlint`, FrontPar: cli.Par{`-q`}, ContainerImage: imgCrieNpm}},

	// noExplicitConstructor and noConstructor unfortunately have problems with CUDA_CALLABLE
	{
		Name:  `cpp`,
		Regex: regexp.MustCompile(`\.cc$|\.cpp$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}, ContainerImage: imgCrieApk},
		Chk: &cli.Lint{
			Bin: `cppcheck`,
			FrontPar: cli.Par{
				`--enable=all`, `--language=c++`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unRegexedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`,
			},
			ContainerImage: imgCrieApk,
		},
	},

	{
		Name:  `cppheaders`,
		Regex: regexp.MustCompile(`\.h$|\.hpp$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}, ContainerImage: imgCrieApk}},

	{
		Name:  `c`,
		Regex: regexp.MustCompile(`\.c$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}, ContainerImage: imgCrieApk},
		Chk: &cli.Lint{
			Bin: `cppcheck`,
			FrontPar: cli.Par{
				`--enable=all`, `--language=c`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unRegexedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`,
			},
			ContainerImage: imgCrieApk,
		},
	},

	{
		Name:  `cmake`,
		Regex: regexp.MustCompile(`CMakeLists.txt$|\.cmake$`),
		Chk:   &cli.Lint{Bin: `cmakelint`, FrontPar: cli.Par{"--config=/home/standards/.config/cmakelintrc"}, ContainerImage: imgCriePip}},

	// TODO Review tools that parse child files - ansiblelint needs to install dependencies similiar to how clang-tidy does
	//{
	//	Name:  `ansible`,
	//	Regex: regexp.MustCompile(`playbook.yml$`),
	//	Chk:   &cli.Lint{Bin: `ansible-lint`, ContainerImage: imgCriePip}},
	//},

	// TODO use v2 with go
	//{
	//	Name:  `dockercompose`,
	//	Regex: regexp.MustCompile(`docker-compose.yml$|docker-compose.yaml$`),
	//	Chk:   &cli.Lint{Bin: `docker-compose`, FrontPar: cli.Par{`-f`}, EndPar: cli.Par{`config`, `-q`}}}
}

func init() {
	//unmarshalProjectSettings(languagesYAML, &Cli.ProjectConfigFile)
	Cli.Crie.Languages = Configuration
}
