package crie

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (s *RunConfiguration) loadFilesGit(args []string) ([]string, error) {
	var outB, errB bytes.Buffer

	gitCmd, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	c := exec.Command(gitCmd, args...)
	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB
	err = c.Run()
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

	for _, reg := range s.IgnoreFiles {
		finallist = RemoveIgnored(finallist, reg.MatchString)
	}

	return finallist, nil
}

// IsRepo checks for a .git folder
func (s *RunConfiguration) IsRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}

func (s *RunConfiguration) fileListRepoChanged() ([]string, error) {
	var outB bytes.Buffer

	gitCmd, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	c := exec.Command(gitCmd,
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

func (s *RunConfiguration) fileListRepoAll() ([]string, error) {
	return s.loadFilesGit([]string{"ls-files"})
}
