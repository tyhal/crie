package crie

import (
	"errors"
	"os"
	"path/filepath"
)

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

	var finalFiles []string
	for _, file := range allFiles {
		if s.Ignore.MatchString(file) {
			finalFiles = append(finalFiles, file)
		}
	}

	return finalFiles, nil
}
