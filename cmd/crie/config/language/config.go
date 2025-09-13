package language

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Languages is the schema for a projects' project file
type Languages struct {
	Languages map[string]Language `json:"languages" yaml:"languages" jsonschema_description:"a map of languages that crie should be able to run"`
}

func merge(src *Languages, dst *Languages) {

	for langName, lang := range src.Languages {
		if _, ok := dst.Languages[langName]; !ok {
			dst.Languages[langName] = lang
		}
	}
}

// NewLanguageConfigFile Creates the project file locally
func NewLanguageConfigFile(path string) error {
	yamlOut, err := yaml.Marshal(Languages{})

	if err != nil {
		return err
	}

	var buf bytes.Buffer
	// TODO add versioning
	buf.WriteString("# yaml-language-server: $schema=https://raw.githubusercontent.com/tyhal/crie/main/res/schema/lang.json\n")
	buf.Write(yamlOut)
	yamlContent := buf.Bytes()
	err = os.WriteFile(path, yamlContent, 0644)

	if err != nil {
		return err
	}

	return nil
}

// LoadFile will attempt to parse a Language schema compatible file and use it to overwrite the builtin defaults
func LoadFile(path string) (*Languages, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &defaultLanguageConfig, nil
	}

	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var c Languages
	if err := yaml.Unmarshal(configData, &c); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	merge(&defaultLanguageConfig, &c)

	return &c, nil
}
