package runner

import (
	"regexp"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter/noop"
)

func Test_getName(t *testing.T) {
	assert.Equal(t, "", getName(nil))
	assert.Equal(t, "noop", getName(&noop.LintNoop{}))
}

func TestRunConfiguration_GetLanguage(t *testing.T) {
	config := &RunConfiguration{
		Languages: map[string]*Language{
			"test": {
				Chk:   &noop.LintNoop{},
				Fmt:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.test$`),
			},
		},
	}

	// Test existing language
	lang, err := config.GetLanguage("test")
	assert.NoError(t, err)
	assert.NotNil(t, lang)

	// Test non-existent language
	_, err = config.GetLanguage("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "language 'nonexistent' not found")
}

func TestRunConfiguration_runLinter(t *testing.T) {
	config := &RunConfiguration{
		Languages: map[string]*Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
		},
	}

	fileList := []string{"test.go"}

	var cleanupGroup sync.WaitGroup
	err := config.runLinter(&cleanupGroup, "go", LintTypeChk, fileList)
	assert.NoError(t, err)

	// Wait for cleanup to complete
	cleanupGroup.Wait()
}

func TestRunConfiguration_runLinters(t *testing.T) {
	tests := []struct {
		name       string
		config     *RunConfiguration
		files      []string
		expectErr  bool
		errMessage string
	}{
		{
			name: "default runLinters - happy path",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:   &noop.LintNoop{},
						Regex: regexp.MustCompile(`\.go$`),
					},
				},
			},
			files:     []string{"test.go"},
			expectErr: false,
		},
		{
			name: "runLinters with single valid language (go)",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:   &noop.LintNoop{},
						Regex: regexp.MustCompile(`\.go$`),
					},
					"test": {
						Chk:   &noop.LintNoop{},
						Regex: regexp.MustCompile(`\.test$`),
					},
				},
				Options: Options{
					Only: "go",
				},
			},
			files:     []string{"test.go"},
			expectErr: false,
		},
		{
			name: "runLinters with nonexistent language in 'Only' option",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:   &noop.LintNoop{},
						Regex: regexp.MustCompile(`\.go$`),
					},
				},
				Options: Options{
					Only: "nonexistent",
				},
			},
			files:      []string{"test.go"},
			expectErr:  true,
			errMessage: "language 'nonexistent' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.runLinters(LintTypeChk, tt.files)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
