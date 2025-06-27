package settings

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

func unmarshalProjectSettings(data []byte, projectSettings *ProjectConfigFile) {
	if err := yaml.Unmarshal(data, projectSettings); err != nil {
		panic(fmt.Sprintf("failed to parse internal language settings: %v", err))
	}
	for i := range projectSettings.Languages {
		lang := &projectSettings.Languages[i]
		lang.Regex = regexp.MustCompile(strings.Join(lang.Match, "|"))
	}
}

// CreateNewProjectSettings Creates the settings file locally
func (cli *CliSettings) CreateNewProjectSettings() {
	yamlOut, err := yaml.Marshal(ProjectConfigFile{})

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(cli.ConfigPath, yamlOut, 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New languages file created: %s\nPlease view this and configure for your repo\n", cli.ConfigPath)
}

// LoadConfigFile load overrides for our projects' settings
func (cli *CliSettings) LoadConfigFile() {
	// TODO
}
