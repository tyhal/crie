package language

import (
	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter/noop"
	"regexp"
	"testing"
)

func TestConfigLanguage_toLanguage(t *testing.T) {
	regex := &Regex{Regexp: regexp.MustCompile(`\.go$`)}
	fmtLinter := Linter{Linter: &noop.LintNoop{}}
	chkLinter := Linter{Linter: &noop.LintNoop{}}

	config := Language{
		Regex: regex,
		Fmt:   fmtLinter,
		Chk:   chkLinter,
	}

	result := config.ToCrieLanguage()

	assert.Equal(t, regex.Regexp, result.Regex)
	assert.Equal(t, fmtLinter.Linter, result.Fmt)
	assert.Equal(t, chkLinter.Linter, result.Chk)
}
