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

func TestConfigLanguage_toLanguageBatch(t *testing.T) {
	regex := &Regex{Regexp: regexp.MustCompile(`\.go$`)}
	tests := []struct {
		name    string
		exec    executor.Instance
		groupBy string
		expErr  string
	}{
		{name: "batch", exec: executor.Instance{Bin: "tool", Batch: true}},
		{name: "batch with chdir", exec: executor.Instance{Bin: "tool", Batch: true, ChDir: true}, expErr: "batch can't be combined with chdir"},
		{name: "batch with no_file_arg", exec: executor.Instance{Bin: "tool", Batch: true, NoFileArg: true}, expErr: "batch can't be combined with chdir"},
		{name: "batch with group_by", exec: executor.Instance{Bin: "tool", Batch: true}, groupBy: "go.mod", expErr: "batch can't be combined with group_by"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Language{
				FileMatch: regex,
				GroupBy:   tt.groupBy,
				Chk:       Linter{Linter: &cli.LintCli{Exec: tt.exec}},
			}
			_, err := config.ToRunFormat()
			if tt.expErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expErr)
			}
		})
	}
}
