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
		par{"rev-list",
			"--no-merges",
			"--count",
			"HEAD"}...)
	c.Stdout = &outB
	if err := c.Run(); err != nil {
		return nil, err
	}

	// Produce string that will  query back all history or only 10 commits
	commitCntStr := strings.Split(outB.String(), "\n")[0]
	commitCnt, err := strconv.Atoi(commitCntStr)
	commitSlice := "HEAD~" + strconv.Itoa(min(commitCnt-1, s.GitDiff)) + "..HEAD"

	args := par{"diff", "--name-only", commitSlice, "."}
	c = exec.Command("git", args...)
	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB
	err = c.Run()

	if err != nil {
		log.WithFields(log.Fields{"type": "stdout"}).Debug(outB)
		log.WithFields(log.Fields{"type": "stderr"}).Debug(errB)
		return nil, err
	}

	return s.loadFileSettings(strings.Split(outB.String(), "\n")), nil
}

func (s *ProjectLintConfiguration) fileListRepoAll() ([]string, error) {
	var outB, errB bytes.Buffer
	args := par{"ls-files"}
	c := exec.Command("git", args...)
	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB
	err := c.Run()

	if err != nil {
		log.WithFields(log.Fields{"type": "stdout"}).Debug(outB)
		log.WithFields(log.Fields{"type": "stderr"}).Debug(errB)
		return nil, err
	}

	return s.loadFileSettings(strings.Split(outB.String(), "\n")), nil
}

func (s *ProjectLintConfiguration) fileListAll() ([]string, error) {

	// Work out where we are
	dir, err := os.Getwd()
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
