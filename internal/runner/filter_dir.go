package runner

import (
	"os"
	"path/filepath"
)

func (s *RunConfiguration) fileListIgnore(allFiles []string) []string {
	if s.Ignore == nil {
		return allFiles
	}
	return Filter(allFiles, false, s.Ignore.MatchString)
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

	return s.fileListIgnore(allFiles), nil
}
