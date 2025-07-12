package settings

import (
	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter/noop"
	"regexp"
	"testing"
)

func TestConfigLanguage_toLanguage(t *testing.T) {
	regex := &ConfigRegex{Regexp: regexp.MustCompile(`\.go$`)}
	fmtLinter := ConfigLinter{Linter: &noop.LintNoop{}}
	chkLinter := ConfigLinter{Linter: &noop.LintNoop{}}

	config := ConfigLanguage{
		Regex: regex,
		Fmt:   fmtLinter,
		Chk:   chkLinter,
	}

	result := config.toLanguage()

	assert.Equal(t, regex.Regexp, result.Regex)
	assert.Equal(t, fmtLinter.Linter, result.Fmt)
	assert.Equal(t, chkLinter.Linter, result.Chk)
}
