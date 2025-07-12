package crie

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
)

// IsRepo checks for a .git folder
func (s *RunConfiguration) IsRepo(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}

func (s *RunConfiguration) fileListRepoChanged(path string) ([]string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}
	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}
	headTree, err := headCommit.Tree()
	if err != nil {
		return nil, err
	}

	splitRef := strings.Split(s.GitTarget, "/")
	if len(splitRef) != 2 {
		return nil, fmt.Errorf("invalid git target must be in form remote/branch")
	}
	targetRef, err := repo.Reference(plumbing.NewRemoteReferenceName(splitRef[0], splitRef[1]), true)
	if err != nil {
		return nil, err
	}
	targetCommit, err := repo.CommitObject(targetRef.Hash())
	if err != nil {
		return nil, err
	}
	targetTree, err := targetCommit.Tree()
	if err != nil {
		return nil, err
	}

	changes, err := object.DiffTree(targetTree, headTree)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, change := range changes {
		if change.To.Name != "" {
			files = append(files, change.To.Name)
		}
	}

	return files, nil
}

func (s *RunConfiguration) fileListRepoAll(path string) ([]string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	idx, err := repo.Storer.Index()
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range idx.Entries {
		files = append(files, entry.Name)
	}

	return files, nil
}
