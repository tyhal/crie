package language

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Config is the entire current Crie Config languages
var Config ConfigLanguages

// ConfigLanguages is the schema for a projects' project file
type ConfigLanguages struct {
	Path      string
	Languages map[string]ConfigLanguage `json:"languages" yaml:"languages" jsonschema_description:"a map of languages that crie should be able to run"`
}

func merge(src *ConfigLanguages, dst *ConfigLanguages) {

	for langName, lang := range src.Languages {
		if _, ok := dst.Languages[langName]; !ok {
			dst.Languages[langName] = lang
		}
	}
}

// NewProjectConfigFile Creates the project file locally
func (c *ConfigLanguages) NewLanguageConfigFile() error {
	yamlOut, err := yaml.Marshal(ConfigLanguages{})

	if err != nil {
		return err
	}

	// TODO output: # yaml-language-server: $schema=./schema.json with the path matching the version of crie being used
	err = os.WriteFile(c.Path, yamlOut, 0644)

	if err != nil {
		return err
	}

	fmt.Printf("New languages file created: %s\nPlease view this and configure for your repo\n", c.Path)
	return nil
}

// LoadFile load overrides for our projects' project
func (c *ConfigLanguages) LoadFile() error {
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		return nil
	}

	configData, err := os.ReadFile(c.Path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", c.Path, err)
	}

	if err := yaml.Unmarshal(configData, c); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", c.Path, err)
	}

	merge(&defaultLanguageConfig, c)

	return nil
}
