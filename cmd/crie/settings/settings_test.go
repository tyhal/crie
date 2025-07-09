package settings

import (
	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/noop"
	"regexp"
	"testing"
)

func TestSaveConfiguration(t *testing.T) {
	fmtLinter := &noop.LintNoop{}
	chkLinter := &noop.LintNoop{}
	regex := regexp.MustCompile(`\.go$`)

	cli := &CliSettings{
		ConfigProject: ConfigProject{
			Ignore: []string{"\\*.tmp"},
			Languages: map[string]ConfigLanguage{
				"go": {
					Regex: &ConfigRegex{Regexp: regex},
					Fmt:   ConfigLinter{Linter: fmtLinter},
					Chk:   ConfigLinter{Linter: chkLinter},
				},
			},
		},
	}

	cli.SaveConfiguration()

	expected := &linter.Language{
		Regex: regex,
		Fmt:   fmtLinter,
		Chk:   chkLinter,
	}

	assert.Equal(t, expected, cli.Crie.Languages["go"])
}
