package api

// TODO Generalise so you can have anything as a 'lint me list'

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *ProjectLintConfiguration) fileListRepoChanged() ([]string, error) {
	var outB, errB bytes.Buffer

	c := exec.Command("git",
		Par{"rev-list",
			"--no-merges",
			"--max-count",
			strconv.Itoa(s.GitDiff),
			"--count",
			"HEAD"}...)
	c.Stdout = &outB
	if err := c.Run(); err != nil {
		return nil, err
	}

	commitSlice := "HEAD~" + strings.Split(outB.String(), "\n")[0] + "..HEAD"

	args := Par{"diff", "--name-only", commitSlice, "."}
	c = exec.Command("git", args...)
	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB
	err := c.Run()

	if err != nil {
		log.WithFields(log.Fields{"type": "stdout"}).Debug(outB.String())
		log.WithFields(log.Fields{"type": "stderr"}).Debug(errB.String())
		return nil, err
	}

	return s.loadFileSettings(strings.Split(outB.String(), "\n")), nil
}

func (s *ProjectLintConfiguration) fileListRepoAll() ([]string, error) {
	var outB, errB bytes.Buffer
	args := Par{"ls-files"}
	c := exec.Command("git", args...)
	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB
	err := c.Run()

	if err != nil {
		log.WithFields(log.Fields{"type": "stdout"}).Debug(outB.String())
		log.WithFields(log.Fields{"type": "stderr"}).Debug(errB.String())
		return nil, err
	}

	return s.loadFileSettings(strings.Split(outB.String(), "\n")), nil
}

func (s *ProjectLintConfiguration) fileListAll() ([]string, error) {

	// Work out where we are
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	allFiles := []string{}

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

	empty, err := isEmpty(".")
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

	return s.loadFileSettings(allFiles), nil
}
