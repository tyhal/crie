package project

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNewProjectConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	cli := &Config{}
	conf := filepath.Join(tempDir, "Config.yml")

	err := cli.NewProjectConfigFile(conf)

	assert.NoError(t, err)
	assert.FileExists(t, conf)
}
