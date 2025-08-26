package language

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigProject_merge(t *testing.T) {
	tests := []struct {
		name     string
		base     ConfigLanguages
		src      ConfigLanguages
		expected ConfigLanguages
	}{
		{
			name:     "add new language",
			base:     ConfigLanguages{Languages: map[string]ConfigLanguage{}},
			src:      ConfigLanguages{Languages: map[string]ConfigLanguage{"go": {}}},
			expected: ConfigLanguages{Languages: map[string]ConfigLanguage{"go": {}}},
		},
		{
			name:     "keep existing language",
			base:     ConfigLanguages{Languages: map[string]ConfigLanguage{"go": {Fmt: ConfigLinter{}}}},
			src:      ConfigLanguages{Languages: map[string]ConfigLanguage{"go": {}}},
			expected: ConfigLanguages{Languages: map[string]ConfigLanguage{"go": {Fmt: ConfigLinter{}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.base.Languages == nil {
				tt.base.Languages = make(map[string]ConfigLanguage)
			}
			merge(&tt.src, &tt.base)

			assert.Equal(t, tt.expected, tt.base)
		})
	}
}
