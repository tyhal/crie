package crie

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	var allFiles []string

	// Create an initial file list
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			allFiles = append(allFiles, path)
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
