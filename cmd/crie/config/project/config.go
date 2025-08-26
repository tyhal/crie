package project

import (
	"fmt"
	language2 "github.com/tyhal/crie/cmd/crie/config/language"
	"gopkg.in/yaml.v3"
	"os"
)

// Config is the entire current Crie Config project
var Config ConfigProject

// ConfigProject are all the things for crie cli
type ConfigProject struct {
	Path    string // 2. user Config overrides defaults
	Quiet   bool
	Verbose bool
	Trace   bool
	JSON    bool
	Ignore  []string `json:"ignore" yaml:"ignore" jsonschema_description:"list of regexes matched against the file list to ignore them (exact paths also work)"`
}

// NewProjectConfigFile Creates the project file locally
func (cli *ConfigProject) NewProjectConfigFile() error {
	yamlOut, err := yaml.Marshal(language2.ConfigLanguages{})

	if err != nil {
		return err
	}

	// TODO output: # yaml-language-server: $schema=./schema.json with the path matching the version of crie being used
	err = os.WriteFile(cli.Path, yamlOut, 0644)

	if err != nil {
		return err
	}

	fmt.Printf("New languages file created: %s\nPlease view this and configure for your repo\n", cli.Path)
	return nil
}

// LoadFile load overrides for our projects' project
func (cli *ConfigProject) LoadFile() error {

	// TODO, do this through viper or something

	if _, err := os.Stat(cli.Path); os.IsNotExist(err) {
		return nil
	}

	configData, err := os.ReadFile(cli.Path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", cli.Path, err)
	}

	if err := yaml.Unmarshal(configData, cli); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", cli.Path, err)
	}

	return nil
}
