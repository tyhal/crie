package project

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (s *LintConfiguration) loadFilesGit(args []string) ([]string, error) {
	var outB, errB bytes.Buffer
	c := exec.Command("git", args...)
	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB
	err := c.Run()
	if err != nil {
		log.WithFields(log.Fields{"type": "stdout"}).Debug(&outB)
		log.WithFields(log.Fields{"type": "stderr"}).Debug(&errB)
		return nil, err
	}

	// Skip files that do not exist at head
	filelist := strings.Split(outB.String(), "\n")
	var finallist []string
	for _, file := range filelist {
		_, err := os.Stat(file)
		if err == nil {
			finallist = append(finallist, file)
		}
	}

	return s.processFilesWithConfig(finallist), nil
}

// IsRepo checks for a .git folder
func (s *LintConfiguration) IsRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}

func (s *LintConfiguration) fileListRepoChanged() ([]string, error) {
	var outB bytes.Buffer

	c := exec.Command("git",
		[]string{"rev-list",
			"--no-merges",
			"--max-count",
			strconv.Itoa(s.GitDiff),
			"--count",
			"HEAD"}...)
	c.Stdout = &outB
	if err := c.Run(); err != nil {
		return nil, err
	}
	commitSlice := "HEAD~" + strings.Split(outB.String(), "\n")[0]
	return s.loadFilesGit([]string{"diff", "--name-only", commitSlice, "."})
}

func (s *LintConfiguration) fileListRepoAll() ([]string, error) {
	return s.loadFilesGit([]string{"ls-files"})
}
