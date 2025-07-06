package settings

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestDefaultProjectConfigYaml(t *testing.T) {
	var config ConfigProject
	err := yaml.Unmarshal(defaultProjectConfigYaml, &config)

	assert.NoError(t, err)
	assert.NotNil(t, config.Languages)
}
