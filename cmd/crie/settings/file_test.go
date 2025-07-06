package settings

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateNewProjectSettings(t *testing.T) {
	tempDir := t.TempDir()
	cli := &CliSettings{ConfigPath: filepath.Join(tempDir, "config.yml")}

	err := cli.CreateNewProjectSettings()

	assert.NoError(t, err)
	assert.FileExists(t, cli.ConfigPath)
}

func TestLoadConfigFile_NoFile(t *testing.T) {
	cli := &CliSettings{ConfigPath: "nonexistent.yml"}

	err := cli.LoadConfigFile()

	assert.NoError(t, err) // Should handle missing file gracefully
}

func TestLoadConfigFile_MergesConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	// Write test config
	testConfig := `ignore: ["*.test"]`
	os.WriteFile(configPath, []byte(testConfig), 0644)

	cli := &CliSettings{
		ConfigPath:    configPath,
		ConfigProject: ConfigProject{Ignore: []string{"*.orig"}},
	}

	err := cli.LoadConfigFile()

	assert.NoError(t, err)
	assert.Contains(t, cli.ConfigProject.Ignore, "*.test")
	assert.Contains(t, cli.ConfigProject.Ignore, "*.orig")
}
