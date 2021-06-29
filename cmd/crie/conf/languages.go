package conf

import (
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/cli"
	proto2 "github.com/tyhal/crie/pkg/linter/proto"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	vale2 "github.com/tyhal/crie/pkg/linter/vale"
	"mvdan.cc/sh/v3/syntax"
	"regexp"
)

// Directory to store default configurations for tools
var confDir = "/etc/crie" // || C:\Program Files\Common Files\crie

var hadolintImg = "docker.io/tyhal/hadolint:0.0.2"
var remarkImg = "docker.io/tyhal/remark:0.0.2"
var terraformImg = "docker.io/hashicorp/terraform:1.0.1"

// LanguageList is a monolithic configuration of all cries standard linters
var LanguageList = []linter.Language{
	{
		Name:  `python`,
		Match: regexp.MustCompile(`\.py$`),
		Fmt:   &cli.Lint{Bin: `autopep8`, FrontPar: cli.Par{`--in-place`, `--aggressive`, `--aggressive`}},
		Chk:   &cli.Lint{Bin: `pylint`, FrontPar: cli.Par{`--rcfile=` + confDir + `/python/.pylintrc`}},
	},
	{
		Name:  `pythondeps`,
		Match: regexp.MustCompile(`requirements.txt$`),
		Fmt:   &cli.Lint{Bin: `pur`, FrontPar: cli.Par{`-r`}},
	},
	{
		Name:  `proto`,
		Match: regexp.MustCompile(`.proto$`),
		Chk:   &proto2.Lint{Fix: false},
		Fmt:   &proto2.Lint{Fix: true},
	},

	// https://github.com/mvdan/sh/releases/download/v1.3.0/shfmt_v1.3.0_linux_amd64
	// https://github.com/koalaman/shellcheck
	{
		Name:  `bash`,
		Match: regexp.MustCompile(`\.bash$`),
		Fmt:   &shfmt.Lint{Language: syntax.LangBash},
		Chk: &cli.Lint{Bin: `shellcheck`, FrontPar: cli.Par{`-x`, `--shell=bash`, `-Calways`},
			Docker: cli.DockerCmd{Image: hadolintImg}}},
	{
		Name:  `sh`,
		Match: regexp.MustCompile(`\.sh$|/script/[^.]*$|^script/[^.]*$`),
		Fmt:   &shfmt.Lint{Language: syntax.LangPOSIX},
		Chk: &cli.Lint{Bin: `shellcheck`, FrontPar: cli.Par{`-x`, `--shell=sh`, `-Calways`},
			Docker: cli.DockerCmd{Image: hadolintImg}}},

	// https://github.com/lukasmartinelli/hadolint
	// TODO use config file when it can be mounted into the docker cmd too
	{
		Name:  `docker`,
		Match: regexp.MustCompile(`Dockerfile$`),
		Chk: &cli.Lint{
			Bin:      `hadolint`,
			FrontPar: cli.Par{`--ignore`, `DL3007`, `--ignore`, `DL3018`, `--ignore`, `DL3016`, `--ignore`, `DL4006`},
			Docker:   cli.DockerCmd{Image: hadolintImg}}},

	//	Fmt:   Lint{`dockfmt`, par{`fmt`, `-w`}, par{}}}

	// https://github.com/adrienverge/yamllint
	{
		Name:  `yml`,
		Match: regexp.MustCompile(`\.yml$|\.yaml$`),
		Chk:   &cli.Lint{Bin: `yamllint`, FrontPar: cli.Par{`-c=` + confDir + `/yaml/.yamllintrc`}}},

	{
		Name:  `terraform`,
		Match: regexp.MustCompile(`\.tf$`),
		Fmt:   &cli.Lint{Bin: `terraform`, Docker: cli.DockerCmd{Image: terraformImg}, FrontPar: cli.Par{`fmt`}},
		Chk:   &cli.Lint{Bin: `terraform`, Docker: cli.DockerCmd{Image: terraformImg}, FrontPar: cli.Par{`fmt`, `-check=true`}}},

	// https://blog.jetbrains.com/webstorm/2017/01/webstorm-2017-1-eap-171-2272/
	// https://github.com/standard/standard
	{
		Name:  `javascript`,
		Match: regexp.MustCompile(`\.js$|\.jsx$`),
		Fmt:   &cli.Lint{Bin: `standard`, FrontPar: cli.Par{`--fix`}},
		Chk:   &cli.Lint{Bin: `standard`}},

	// https://golang.org/cmd/gofmt/
	{
		Name:  `golang`,
		Match: regexp.MustCompile(`\.go$`),
		Fmt:   &cli.Lint{Bin: `gofmt`, FrontPar: cli.Par{`-l`, `-w`}},
		Chk:   &cli.Lint{Bin: `golint`, FrontPar: cli.Par{`-set_exit_status`}},
	},

	// https://github.com/wooorm/remark-lint
	{
		Name:  `markdown`,
		Match: regexp.MustCompile(`\.md$`),
		Fmt:   &cli.Lint{Bin: `remark`, FrontPar: cli.Par{`--use`, `remark-preset-lint-recommended`}, EndPar: cli.Par{`-o`}, Docker: cli.DockerCmd{Image: remarkImg}},
		Chk:   vale2.NewValeLint(confDir + `/markdown/.vale.ini`)},

	{
		Name:  `asciidoctor`,
		Match: regexp.MustCompile(`\.adoc$`),
		Chk:   vale2.NewValeLint(confDir + `/markdown/.vale.ini`)},

	// https://github.com/zaach/jsonlint
	{
		Name:  `json`,
		Match: regexp.MustCompile(`\.json$|\.JSON$`),
		Fmt:   &cli.Lint{Bin: `jsonlint`, FrontPar: cli.Par{`-i`, `-s`, `-c`, `-q`}},
		Chk:   &cli.Lint{Bin: `jsonlint`, FrontPar: cli.Par{`-q`}}},

	// noExplicitConstructor and noConstructor unfortunately have problems with CUDA_CALLABLE
	{
		Name:  `cpp`,
		Match: regexp.MustCompile(`\.cc$|\.cpp$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}},
		Chk:   &cli.Lint{Bin: `cppcheck`, FrontPar: cli.Par{`--enable=all`, `--language=c++`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}}},

	{
		Name:  `cppheaders`,
		Match: regexp.MustCompile(`\.h$|\.hpp$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}}},

	{
		Name:  `c`,
		Match: regexp.MustCompile(`\.c$`),
		Fmt:   &cli.Lint{Bin: `clang-format`, FrontPar: cli.Par{`-style=file`, `-i`}},
		Chk:   &cli.Lint{Bin: `cppcheck`, FrontPar: cli.Par{`--enable=all`, `--language=c`, `--suppress=operatorEqRetRefThis`, `--suppress=operatorEq`, `--suppress=noExplicitConstructor`, `--suppress=unmatchedSuppression`, `--suppress=missingInclude`, `--suppress=unusedFunction`, `--suppress=noConstructor`, `--inline-suppr`, `--error-exitcode=1`}}},

	{
		Name:  `cmake`,
		Match: regexp.MustCompile(`CMakeLists.txt$|\.cmake$`),
		Chk:   &cli.Lint{Bin: `cmakelint`, FrontPar: cli.Par{`--config=` + confDir + `/cmake/.cmakelintrc`}}},

	// TODO Review tools that parse child files - ansiblelint needs to install dependencies similiar to how clang-tidy does
	{
		Name:  `ansible`,
		Match: regexp.MustCompile(`playbook.yml$`),
		//Chk:   &imp.Lint{Bin: `ansible-lint`}
	},

	{
		Name:  `dockercompose`,
		Match: regexp.MustCompile(`docker-compose.yml$|docker-compose.yaml$`),
		Chk:   &cli.Lint{Bin: `docker-compose`, FrontPar: cli.Par{`-f`}, EndPar: cli.Par{`config`, `-q`}}}}
