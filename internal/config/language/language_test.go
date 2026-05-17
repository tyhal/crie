package language

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/cli/executor"

	"github.com/tyhal/crie/pkg/cli"
	"github.com/tyhal/crie/pkg/noop"
)

func TestConfigLanguage_toLanguage(t *testing.T) {
	regex := &Regex{Regexp: regexp.MustCompile(`\.go$`)}
	fmtLinter := Linter{Linter: &cli.LintCli{Exec: executor.Instance{WillWrite: false}}}
	chkLinter := Linter{Linter: &noop.LintNoop{}}

	config := Language{
		FileMatch: regex,
		Fmt:       fmtLinter,
		Chk:       chkLinter,
	}

	result, err := config.ToRunFormat()
	assert.NoError(t, err)
	assert.Equal(t, regex.Regexp, result.FileMatch)
	assert.Equal(t, chkLinter.Linter, result.Chk)

	if assert.IsType(t, &cli.LintCli{}, result.Fmt) {
		assert.True(t, result.Fmt.(*cli.LintCli).Exec.WillWrite, "expected WillWrite to be enabled by default for formatters")
	}
}
