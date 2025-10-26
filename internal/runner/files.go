package runner

import (
	"errors"
	"os"

	"github.com/tyhal/crie/internal/filelist"
)

func (s *RunConfiguration) filterIgnoredFiles(files []string) []string {
	if s.Ignore == nil {
		return files
	}
	return Filter(files, false, s.Ignore.MatchString)
}

func (s *RunConfiguration) getAllFilesList() ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.New("could not get current working directory")
	}

	if repo, err := filelist.GetGitRepo(dir); err == nil && repo != nil {
		return filelist.FromGitRepo(dir, repo, s.Options.GitDiff, s.Options.GitTarget)
	}

	// Check if the user asked for git diffs when not in a repo
	if s.Options.GitDiff {
		return nil, errors.New("you do not appear to be in a git repository")
	}
	return filelist.FromDir(dir)

}

// GetFileList returns a list of files based on
func (s *RunConfiguration) getFileList() ([]string, error) {
	allFiles, err := s.getAllFilesList()
	if err != nil {
		return nil, err
	}
	return s.filterIgnoredFiles(allFiles), nil
}
