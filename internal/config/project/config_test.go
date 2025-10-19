package project

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProjectConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	cli := &Config{}
	conf := filepath.Join(tempDir, "Config.yml")

	err := cli.NewProjectConfigFile(conf)

	assert.NoError(t, err)
	assert.FileExists(t, conf)
}
