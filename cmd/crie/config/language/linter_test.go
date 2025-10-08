package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfigLinter_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name       string
		yaml       string
		expectName string
		wantErr    bool
	}{
		{
			name: "cli linter",
			yaml: `type: cli
exec: 
  bin: "eslint"`,
			expectName: "eslint",
			wantErr:    false,
		},
		{
			name: "shfmt linter",
			yaml: `type: shfmt
language: bash`,
			expectName: "shfmt",
			wantErr:    false,
		},
		{
			name:       "noop linter",
			yaml:       `type: noop`,
			expectName: "noop",
			wantErr:    false,
		},
		{
			name:    "unknown type",
			yaml:    `type: unknown`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cl Linter
			err := yaml.Unmarshal([]byte(tt.yaml), &cl)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectName, cl.Linter.Name())
			}
		})
	}
}
