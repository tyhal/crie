package settings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigProjectMerge(t *testing.T) {
	tests := []struct {
		name     string
		base     ConfigProject
		src      ConfigProject
		expected ConfigProject
	}{
		{
			name:     "add new language",
			base:     ConfigProject{Languages: map[string]ConfigLanguage{}},
			src:      ConfigProject{Languages: map[string]ConfigLanguage{"go": {}}},
			expected: ConfigProject{Languages: map[string]ConfigLanguage{"go": {}}},
		},
		{
			name:     "keep existing language",
			base:     ConfigProject{Languages: map[string]ConfigLanguage{"go": {Fmt: ConfigLinter{}}}},
			src:      ConfigProject{Languages: map[string]ConfigLanguage{"go": {}}},
			expected: ConfigProject{Languages: map[string]ConfigLanguage{"go": {Fmt: ConfigLinter{}}}},
		},
		{
			name:     "merge ignore lists",
			base:     ConfigProject{Ignore: []string{"*.tmp"}},
			src:      ConfigProject{Ignore: []string{"*.log"}},
			expected: ConfigProject{Languages: map[string]ConfigLanguage{}, Ignore: []string{"*.tmp", "*.log"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.base.Languages == nil {
				tt.base.Languages = make(map[string]ConfigLanguage)
			}
			tt.base.merge(tt.src)

			assert.Equal(t, tt.expected, tt.base)
		})
	}
}
