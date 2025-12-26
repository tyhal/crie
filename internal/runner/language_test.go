package runner

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter/noop"
)

func TestLanguages_Show(t *testing.T) {
	tests := []struct {
		name      string
		languages Languages
		expected  []string // Substrings we expect to be in the output
	}{
		{
			name: "two languages",
			languages: Languages{
				"python": {
					FileMatch: regexp.MustCompile(`\.py$`),
					Chk:       &noop.LintNoop{},
				},
				"go": {
					FileMatch: regexp.MustCompile(`\.go$`),
					Fmt:       &noop.LintNoop{},
				},
			},
			expected: []string{"go", "noop", "noop", "\\.go$", "python", "noop", "noop", "\\.py$"},
		},
		{
			name: "one language",
			languages: Languages{
				"yaml": {
					FileMatch: regexp.MustCompile(`\.ya?ml$`),
					Chk:       &noop.LintNoop{},
				},
			},
			expected: []string{"yaml", "noop", "", "\\.ya?ml$"},
		},
		{
			name:      "no languages",
			languages: Languages{},
			expected:  []string{"LANGUAGE", "CHECKER", "FORMATTER", "ASSOCIATED FILES"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := tt.languages.Show(&buf)
			assert.NoError(t, err)

			output := buf.String()

			for _, expectedPart := range tt.expected {
				assert.Contains(t, output, expectedPart)
			}
		})
	}
}

func Test_getName(t *testing.T) {
	assert.Empty(t, getName(nil))
	assert.Equal(t, "noop", getName(&noop.LintNoop{}))
}
