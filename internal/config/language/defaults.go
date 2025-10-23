package language

import (
	_ "embed" // Used to embed our default configurations
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed defaults.yml
var defaultLanguageConfigYaml []byte
var defaultLanguageConfig Languages

func init() {
	err := yaml.Unmarshal(defaultLanguageConfigYaml, &defaultLanguageConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to parse internal language project: %v", err))
	}
}
