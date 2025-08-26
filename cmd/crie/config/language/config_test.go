package language

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigProject_merge(t *testing.T) {
	tests := []struct {
		name     string
		base     Languages
		src      Languages
		expected Languages
	}{
		{
			name:     "add new language",
			base:     Languages{Languages: map[string]Language{}},
			src:      Languages{Languages: map[string]Language{"go": {}}},
			expected: Languages{Languages: map[string]Language{"go": {}}},
		},
		{
			name:     "keep existing language",
			base:     Languages{Languages: map[string]Language{"go": {Fmt: Linter{}}}},
			src:      Languages{Languages: map[string]Language{"go": {}}},
			expected: Languages{Languages: map[string]Language{"go": {Fmt: Linter{}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.base.Languages == nil {
				tt.base.Languages = make(map[string]Language)
			}
			merge(&tt.src, &tt.base)

			assert.Equal(t, tt.expected, tt.base)
		})
	}
}
