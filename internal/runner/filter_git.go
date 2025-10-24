package runner

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/tyhal/crie/pkg/errchain"
)

// IsRepo checks for a .git folder
func (s *RunConfiguration) IsRepo(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}

func getReference(target string) (plumbing.ReferenceName, error) {
	splitRef := strings.Split(target, "/")
	switch len(splitRef) {
	case 1:
		return plumbing.NewBranchReferenceName(splitRef[0]), nil
	case 2:
		return plumbing.NewRemoteReferenceName(splitRef[0], splitRef[1]), nil
	default:
		return "", fmt.Errorf("invalid git target must be in form 'remote/branch' or 'branch'")
	}
}

func (s *RunConfiguration) changesToTarget(repo *git.Repository, target string) ([]string, error) {
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
	refName, err := getReference(target)
	if err != nil {
		return nil, err
	}
	targetRef, err := repo.Reference(refName, true)
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

func (s *RunConfiguration) changesUncommitted(repo *git.Repository) ([]string, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}
	var files []string
	for filePath, fileStatus := range status {
		if fileStatus.Staging != git.Unmodified || fileStatus.Worktree != git.Unmodified {
			files = append(files, filePath)
		}
	}
	return files, nil
}

func (s *RunConfiguration) fileListRepoChanged(path string) ([]string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	var files []string
	if s.Options.GitTarget == "" {
		files, err = s.changesUncommitted(repo)
		if err != nil {
			return nil, errchain.From(err).Error("getting uncommitted changes")
		}
	} else {
		files, err = s.changesToTarget(repo, s.Options.GitTarget)
		if err != nil {
			return nil, errchain.From(err).ErrorF("getting changes to target %s", s.Options.GitTarget)
		}
	}

	return s.fileListIgnore(files), nil
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

	return s.fileListIgnore(files), nil
}
