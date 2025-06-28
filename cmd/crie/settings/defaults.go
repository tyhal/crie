package settings

import (
	_ "embed" // Used to embed our default configurations
	"fmt"
	"gopkg.in/yaml.v3"
)

//go:embed defaults.yml
var defaultProjectConfigYaml []byte

func init() {
	err := yaml.Unmarshal(defaultProjectConfigYaml, &Cli.ConfigProject)
	if err != nil {
		panic(fmt.Sprintf("failed to parse internal language settings: %v", err))
	}
}
