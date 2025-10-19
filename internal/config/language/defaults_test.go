package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDefaultProjectConfigYaml(t *testing.T) {
	var config Languages
	err := yaml.Unmarshal(defaultLanguageConfigYaml, &config)

	assert.NoError(t, err)
	assert.NotNil(t, config.Languages)
}
