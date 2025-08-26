package project

import (
	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/cmd/crie/config/language"
	"os"
	"path/filepath"
	"testing"
)

// TODO move to project Config test

func TestCreateNewProjectSettings(t *testing.T) {
	tempDir := t.TempDir()
	cli := &Config{Path: filepath.Join(tempDir, "Config.yml")}

	err := cli.NewProjectConfigFile()

	assert.NoError(t, err)
	assert.FileExists(t, cli.Path)
}

func TestLoadConfigFile_NoFile(t *testing.T) {
	cli := &Config{Path: "nonexistent.yml"}

	err := cli.LoadFile()

	assert.NoError(t, err) // Should handle missing file gracefully
}

func TestLoadConfigFile_MergesConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "Config.yml")

	// Write test Config
	testConfig := `ignore: ["*.test"]`
	os.WriteFile(configPath, []byte(testConfig), 0644)

	cli := &Config{
		Path:          configPath,
		ConfigProject: language.Languages{Ignore: []string{"*.orig"}},
	}

	err := cli.LoadFile()

	assert.NoError(t, err)
	assert.Contains(t, cli.ConfigProject.Ignore, "*.test")
	assert.Contains(t, cli.ConfigProject.Ignore, "*.orig")
}
