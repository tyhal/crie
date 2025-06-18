package crie

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

// loadFileList returns all valid files that have also been filtered by the config
func (s *RunConfiguration) loadFileList() {

	var fileList []string
	var err error

	if s.IsRepo() {

		if s.GitDiff > 0 {
			// Get files changed in last s.GitDiff commits
			fileList, err = s.fileListRepoChanged()
		} else {
			// Get all files in git repo
			fileList, err = s.fileListRepoAll()
		}
	} else {

		// Check if the user asked for git diffs when not in a repo
		if s.GitDiff > 0 {
			log.Fatal("This is not a git repo you are in")
		}

		// Generic grab all the files
		fileList, err = s.fileListAll()
	}
	if err != nil {
		log.Fatal(err.Error())
	} else {
		s.fileList = fileList
	}
}

// CreateNewProjectSettings Creates the settings file locally
func CreateNewProjectSettings(confpath string) {
	yamlOut, err := yaml.Marshal(ProjectSettings{})

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(confpath, yamlOut, 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New languages file created: %s\nPlease view this and configure for your repo\n", confpath)
}
