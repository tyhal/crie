package project

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/internal/helper"
	"github.com/tyhal/crie/internal/settings"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
)

func (s *LintConfiguration) processFilesWithConfig(files []string) []string {
	f, err := os.Open(s.ConfPath)

	if err != nil {
		log.Fatal(err)
	}

	m := settings.FileSettings{}

	err = yaml.NewDecoder(f).Decode(&m)

	if err != nil {
		log.Fatal("Failed to parse (" + s.ConfPath + "): " + err.Error())
	}

	for _, ignReg := range m.Ignore {
		reg, err := regexp.Compile(ignReg)

		if err != nil {
			log.Fatal(err)
		}

		files = helper.RemoveIgnored(files, reg.MatchString)
	}

	return files
}

func (s *LintConfiguration) fileListAll() ([]string, error) {

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

	empty, err := helper.IsEmpty(".")
	if err != nil {
		return nil, err
	}

	if empty {
		return nil, errors.New("this is an empty folder")
	}

	// If there is a config then parse the files through it
	if _, err := os.Stat(s.ConfPath); err != nil {
		return allFiles, nil
	}

	return s.processFilesWithConfig(allFiles), nil
}
