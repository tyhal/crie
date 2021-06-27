package project

import (
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/internal/settings"
	"github.com/tyhal/crie/pkg/crie/linter"
	"os"
)

// LintConfiguration is what is required for an entire project to be linted
type LintConfiguration struct {
	ConfPath        string
	lintType        string
	ContinueOnError bool
	ShowPasses      bool
	Languages       []linter.Language
	GitDiff         int
	SingleLang      string
	fileList        []string
}

// loadFileList returns all valid files that have also been filtered by the config
func (s *LintConfiguration) loadFileList() {

	var fileList []string
	var err error

	if s.IsRepo() {
		// If we are a repo without a configuration then force it upon the project
		if _, err := os.Stat(s.ConfPath); err != nil {
			settings.CreateNewFileSettings(s.ConfPath)
		}

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
