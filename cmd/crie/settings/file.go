package settings

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// CreateNewProjectSettings Creates the settings file locally
func (cli *CliSettings) CreateNewProjectSettings() error {
	yamlOut, err := yaml.Marshal(ConfigProject{})

	if err != nil {
		return err
	}

	// TODO output: # yaml-language-server: $schema=./schema.json with the path matching the version of crie being used
	err = os.WriteFile(cli.ConfigPath, yamlOut, 0644)

	if err != nil {
		return err
	}

	fmt.Printf("New languages file created: %s\nPlease view this and configure for your repo\n", cli.ConfigPath)
	return nil
}

// LoadConfigFile load overrides for our projects' settings
func (cli *CliSettings) LoadConfigFile() error {
	if _, err := os.Stat(cli.ConfigPath); os.IsNotExist(err) {
		return nil
	}

	configData, err := os.ReadFile(cli.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", cli.ConfigPath, err)
	}

	var userConfig ConfigProject
	if err := yaml.Unmarshal(configData, &userConfig); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", cli.ConfigPath, err)
	}

	cli.ConfigProject.merge(userConfig)

	return nil
}
