package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

func (s *ProjectLintConfiguration) loadFileSettings(files []string) []string {
	f, err := os.Open(s.ConfPath)

	if err != nil {
		log.Fatal(err)
	}

	m := fileSettings{}

	err = yaml.NewDecoder(f).Decode(&m)

	if err != nil {
		log.Fatal("Failed to parse (" + s.ConfPath + "): " + err.Error())
	}

	for _, ignReg := range m.Ignore {
		reg, err := regexp.Compile(ignReg)

		if err != nil {
			log.Fatal(err)
		}

		files = removeIgnored(files, reg.MatchString)
	}

	// Add more project roots
	projDirs = m.ProjDirs
	projDirs = append(projDirs, ".")

	return files
}

func createFileSettings(confpath string) {
	yamlOut, err := yaml.Marshal(fileSettings{})

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(confpath, yamlOut, 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New languages fileSettings created: %s\nPlease view this and configure for your repo\n", confpath)
}
