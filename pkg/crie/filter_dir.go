package crie

import (
	"errors"
	"os"
	"path/filepath"
)

// RemoveIgnored Narrows down the list by returning only results that do not match the match in the settings file
func RemoveIgnored(list []string, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range list {
		result := f(entry)
		_, err := os.Stat(entry)
		if !result && err == nil {
			filteredLists = append(filteredLists, entry)
		}
	}
	return filteredLists
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

	for _, reg := range s.IgnoreFiles {
		allFiles = RemoveIgnored(allFiles, reg.MatchString)
	}

	return allFiles, nil
}
