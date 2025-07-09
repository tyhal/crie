package crie

import (
	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/noop"
	"regexp"
	"sync"
	"testing"
)

func TestRunConfiguration_GetLanguage(t *testing.T) {
	config := &RunConfiguration{
		Languages: map[string]*linter.Language{
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
	assert.Contains(t, err.Error(), "language not found")
}

func TestRunConfiguration_runLinter(t *testing.T) {
	config := &RunConfiguration{
		lintType: "chk",
		fileList: []string{"test.go"},
		Languages: map[string]*linter.Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
		},
	}

	var cleanupGroup sync.WaitGroup
	err := config.runLinter(&cleanupGroup, "go", config.Languages["go"])
	assert.NoError(t, err)

	// Wait for cleanup to complete
	cleanupGroup.Wait()
}

func TestRunConfiguration_Run(t *testing.T) {
	config := &RunConfiguration{
		Languages: map[string]*linter.Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
		},
		fileList: []string{"test.go"},
	}

	err := config.Run("chk")
	assert.NoError(t, err)
}

func TestRunConfiguration_Run_WithSingleLang(t *testing.T) {
	// TODO Run should produce a standard report possibly aligned with reporting logs
	// with the report we can compare what was actually run
	config := &RunConfiguration{
		Languages: map[string]*linter.Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
			"test": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.test$`),
			},
		},
		SingleLang: "go",
		fileList:   []string{"test.go"},
	}

	err := config.Run("chk")
	assert.NoError(t, err)
}

func TestRunConfiguration_Run_NonexistentSingleLang(t *testing.T) {
	config := &RunConfiguration{
		Languages: map[string]*linter.Language{
			"go": {
				Chk:   &noop.LintNoop{},
				Regex: regexp.MustCompile(`\.go$`),
			},
		},
		SingleLang: "nonexistent",
		fileList:   []string{"test.go"},
	}

	err := config.Run("chk")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "language not found")
}
