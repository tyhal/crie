package filelist

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetGitRepo opens a git repository rooted at dir (or its parents when .git is detected)
// and returns the repository handle or a wrapped error with context.
func GetGitRepo(dir string) (*git.Repository, error) {
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("opening git repo: %w", err)
	}
	return repo, nil
}

// FromGitRepo returns a list of repository files relative to dir.
// If diff is true and target is empty, it returns uncommitted changes.
// If both diff and target are set, it returns changes compared to the target reference.
// If neither is set, it returns all tracked files. An error is returned for invalid combinations.
func FromGitRepo(dir string, repo *git.Repository, diff bool, target string) ([]string, error) {
	if !path.IsAbs(dir) {
		return nil, errors.New("dir must be an absolute path inside the git repo")
	}

	var files []string
	var err error
	action := "determining file list from git"
	targetSet := target != ""

	switch {
	case diff && targetSet:
		action = fmt.Sprintf("getting changes to target '%s'", target)
		files, err = changesFromTarget(repo, target)
	case diff:
		action = "getting uncommitted changes"
		files, err = changesUncommitted(repo)
	case targetSet:
		err = errors.New("invalid parameter combination with target set and diff unset")
	default:
		action = "getting all files"
		files, err = all(repo)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", action, err)
	}

	prefix, err := getPrefixInRepo(dir, repo)
	return prefixFilter(prefix, files...), nil
}

func getRepoPath(repo *git.Repository) (string, error) {
	wt, err := repo.Worktree()
	if err == nil {
		return wt.Filesystem.Root(), nil
	}
	return "", fmt.Errorf("cannot determine repo path")
}

func relInside(base, curr string) (string, error) {
	rel, err := filepath.Rel(base, curr)
	if err != nil {
		return "", err
	}
	if rel != "." && strings.HasPrefix(rel, "..") {
		return "", errors.New("path is not inside the base path")
	}
	return rel, nil
}

func getPrefixInRepo(dir string, repo *git.Repository) (string, error) {
	repoPath, err := getRepoPath(repo)
	if err != nil {
		return "", err
	}
	prefix, err := relInside(repoPath, dir)
	if err != nil {
		return "", err
	}
	return prefix, nil
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

func prefixFilter(prefix string, filepaths ...string) []string {
	if prefix == "." {
		return filepaths
	}
	// to ensure files don't have a leading slash
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	var filtered []string
	for _, file := range filepaths {
		if strings.HasPrefix(file, prefix) {
			filtered = append(filtered, strings.TrimPrefix(file, prefix))
		}
	}
	return filtered
}

func changesFromTarget(repo *git.Repository, target string) ([]string, error) {
	if repo == nil {
		return nil, errors.New("no repo provided")
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

	files := make([]string, 0, len(changes))
	for _, change := range changes {
		if change.To.Name != "" {
			files = append(files, change.To.Name)
		}
	}
	return files, nil
}

func changesUncommitted(repo *git.Repository) ([]string, error) {
	if repo == nil {
		return nil, errors.New("no repo provided")
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(status))
	for filePath, fileStatus := range status {
		if fileStatus.Staging != git.Unmodified || fileStatus.Worktree != git.Unmodified {
			files = append(files, filePath)
		}
	}
	return files, nil
}

// all returns all files in the repo, dir is an absolute path to the working directory
func all(repo *git.Repository) ([]string, error) {
	if repo == nil {
		return nil, errors.New("no repo provided")
	}

	idx, err := repo.Storer.Index()
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(idx.Entries))
	for _, e := range idx.Entries {
		files = append(files, e.Name)
	}
	return files, nil
}
