package language

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestDefaultProjectConfigYaml(t *testing.T) {
	var config ConfigLanguages
	err := yaml.Unmarshal(defaultLanguageConfigYaml, &config)

	assert.NoError(t, err)
	assert.NotNil(t, config.Languages)
}
