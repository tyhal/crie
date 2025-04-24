package settings

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

// FileSettings simply adjusts what we include in a normal lint
type FileSettings struct {
	Ignore   []string `yaml:"ignore"`
	ProjDirs []string `yaml:"proj_dirs"`
}

// CreateNewFileSettings CreateNewFileSettings
func CreateNewFileSettings(confpath string) {
	yamlOut, err := yaml.Marshal(FileSettings{})

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(confpath, yamlOut, 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New languages file created: %s\nPlease view this and configure for your repo\n", confpath)
}
