package runner

import (
	"errors"
	"fmt"
	"regexp"
)

// Options are the core flags and settings to change execution
type Options struct {
	Continue      bool   `json:"continue" yaml:"continue"`
	Passes        bool   `json:"passes" yaml:"passes"`
	GitTarget     string `json:"gitTarget" yaml:"gitTarget"`
	GitDiff       bool   `json:"gitDiff" yaml:"gitDiff"`
	Only          string `json:"only" yaml:"only"`
	StrictLogging bool   `json:"-" yaml:"-"`
}

// RunConfiguration is the entire working set of information to process a project
type RunConfiguration struct {
	Options   Options
	Ignore    *regexp.Regexp
	Languages Languages
}

// getFileList returns all valid files that have also been filtered by the project
func (s *RunConfiguration) getFileList() (fileList []string, err error) {
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
			return nil, errors.New("you do not appear to be in a git repository")
		}

		// Generic grab all the files
		fileList, err = s.fileListAll()
	}
	return
}

// GetLanguage returns the Language configuration by its name or an error if it does not exist.
func (s *RunConfiguration) GetLanguage(only string) (*Language, error) {
	if lang, ok := s.Languages[only]; ok {
		return lang, nil
	}
	return nil, fmt.Errorf("language '%s' not found", only)
}
