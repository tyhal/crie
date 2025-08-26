package language

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestConfigRegex_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected string
		wantErr  bool
	}{
		{
			name:     "single pattern",
			yaml:     `["\\.go$"]`,
			expected: `\.go$`,
		},
		{
			name:     "multiple patterns",
			yaml:     `["\\.go$", "\\.js$"]`,
			expected: `\.go$|\.js$`,
		},
		{
			name: "empty array",
			yaml: `[]`,
		},
		{
			name:    "invalid yaml",
			yaml:    `invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cr ConfigRegex
			err := yaml.Unmarshal([]byte(tt.yaml), &cr)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.expected != "" {
				assert.Equal(t, tt.expected, cr.Regexp.String())
			} else {
				assert.Nil(t, cr.Regexp)
			}
		})
	}
}
