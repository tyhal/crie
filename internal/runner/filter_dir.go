package runner

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func (s *RunConfiguration) fileListIgnore(allFiles []string) []string {
	if s.Ignore == nil {
		return allFiles
	}

	var filteredFiles []string

	for _, file := range allFiles {
		if !s.Ignore.MatchString(file) {
			filteredFiles = append(filteredFiles, file)
		} else {
			log.Debugf("- ignoring file %s", file)
		}
	}

	return filteredFiles
}

func (s *RunConfiguration) fileListAll() ([]string, error) {

	// Work out where we are
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var allFiles []string

	// Create an initial file list
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			allFiles = append(allFiles, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	empty, err := IsEmpty(".")
	if err != nil {
		return nil, err
	}

	if empty {
		return nil, errors.New("this is an empty folder")
	}

	return s.fileListIgnore(allFiles), nil
}
