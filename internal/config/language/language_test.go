package language

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter/noop"
)

func TestConfigLanguage_toLanguage(t *testing.T) {
	regex := &Regex{Regexp: regexp.MustCompile(`\.go$`)}
	fmtLinter := Linter{Linter: &noop.LintNoop{}}
	chkLinter := Linter{Linter: &noop.LintNoop{}}

	config := Language{
		FileMatch: regex,
		Fmt:       fmtLinter,
		Chk:       chkLinter,
	}

	result, err := config.ToCrieLanguage()
	assert.NoError(t, err)
	assert.Equal(t, regex.Regexp, result.FileMatch)
	assert.Equal(t, fmtLinter.Linter, result.Fmt)
	assert.Equal(t, chkLinter.Linter, result.Chk)
}
