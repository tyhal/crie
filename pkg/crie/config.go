package crie

import (
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/crie/linter"
	"regexp"
)

// Options are the core flags and settings to change execution
type Options struct {
	Continue      bool   `json:"continue" yaml:"continue"`
	Passes        bool   `json:"passes" yaml:"passes"`
	GitTarget     string `json:"gitTarget" yaml:"gitTarget"`
	GitDiff       bool   `json:"gitDiff" yaml:"gitDiff"`
	Lang          string `json:"lang" yaml:"lang"`
	StrictLogging bool   `json:"-" yaml:"-"`
}

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	Options   Options
	Ignore    *regexp.Regexp
	Languages Languages
	lintType  string
	fileList  []string
}

// Languages store the name to a singular language configuration within crie
type Languages map[string]*linter.Language

// loadFileList returns all valid files that have also been filtered by the project
func (s *RunConfiguration) loadFileList() {

	var fileList []string
	var err error

	if s.IsRepo(".") {
		if s.Options.GitDiff {
			// Get files changed in last s.GitDiff commits
			fileList, err = s.fileListRepoChanged(".")
		} else {
			// Get all files in git repo
			fileList, err = s.fileListRepoAll(".")
		}
	} else {

		// Check if the user asked for git diffs when not in a repo
		if s.Options.GitDiff {
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
