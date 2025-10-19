package runner

import (
	"regexp"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter/noop"
)

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
		fileList: []string{"test.go"},
		Languages: map[string]*Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
		},
	}

	var cleanupGroup sync.WaitGroup
	err := config.runLinter(&cleanupGroup, "go", LintTypeChk)
	assert.NoError(t, err)

	// Wait for cleanup to complete
	cleanupGroup.Wait()
}

func TestRunConfiguration_Run(t *testing.T) {
	config := &RunConfiguration{
		Languages: map[string]*Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
		},
		fileList: []string{"test.go"},
	}

	err := config.Run(LintTypeChk)
	assert.NoError(t, err)
}

func TestRunConfiguration_Run_WithSingleLang(t *testing.T) {
	// TODO Run should produce a standard report possibly aligned with reporting logs
	// with the report we can compare what was actually run
	config := &RunConfiguration{
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
		fileList: []string{"test.go"},
	}

	err := config.Run(LintTypeChk)
	assert.NoError(t, err)
}

// TODO table test
func TestRunConfiguration_Run_NonexistentSingleLang(t *testing.T) {
	config := &RunConfiguration{
		Languages: map[string]*Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
		},
		Options: Options{
			Only: "nonexistent",
		},
		fileList: []string{"test.go"},
	}

	err := config.Run(LintTypeChk)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "language 'nonexistent' not found")
}
