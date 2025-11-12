// Package folding provides a way to fold / collapse log messages when it detects it's running in a CI environment'
package folding

import (
	"io"
	"os"
)

// Folder represents a CI environment that supports folding / collapsing log messages
type Folder interface {
	Start(file string, msg string, open bool) (string, error)
	Stop(id string) error
}

func isEnvSet(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

func newFrom(w io.Writer, isSet func(string) bool) Folder {
	switch {
	case isSet("GITHUB_ACTIONS"):
		return NewGithub(w)
	case isSet("GITLAB_CI"):
		return NewGitlab(w)
	default:
		return NewPlain(w)
	}
}

// NewW returns a Folding instance based on a custom writer
func NewW(w io.Writer) Folder {
	return newFrom(w, isEnvSet)
}

// New returns a Folding instance based on if a CI environment variable is set
func New() Folder {
	return newFrom(os.Stdout, isEnvSet)
}
