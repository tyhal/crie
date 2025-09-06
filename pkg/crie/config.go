package crie

import (
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/crie/linter"
	"regexp"
)

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	lintType        string
	ContinueOnError bool
	StrictLogging   bool
	ShowPasses      bool
	Ignore          *regexp.Regexp
	Languages       Languages
	GitDiff         bool
	GitTarget       string
	SingleLang      string
	fileList        []string
}

// Languages store the name to a singular language configuration within crie
type Languages map[string]*linter.Language

// loadFileList returns all valid files that have also been filtered by the project
func (s *RunConfiguration) loadFileList() {

	var fileList []string
	var err error

	if s.IsRepo(".") {
		if s.GitDiff {
			// Get files changed in last s.GitDiff commits
			fileList, err = s.fileListRepoChanged(".")
		} else {
			// Get all files in git repo
			fileList, err = s.fileListRepoAll(".")
		}
	} else {

		// Check if the user asked for git diffs when not in a repo
		if s.GitDiff {
			log.Fatal("You do not appear to be in a git repository")
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
