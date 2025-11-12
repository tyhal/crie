package folding

import (
	"fmt"
	"io"
)

type githubFolder struct {
	io.Writer
}

// NewGithub uses the GitHub Actions log syntax
func NewGithub(w io.Writer) Folder {
	return &githubFolder{w}
}

func (g githubFolder) Start(file string, msg string, _ bool) (string, error) {
	_, err := fmt.Fprintf(g, "::error file=%s::%s\n", file, file)
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprintf(g, "::group::%s see logs\n", msg)
	return "", err
}

func (g githubFolder) Stop(_ string) error {
	_, err := fmt.Fprintf(g, "::endgroup::\n")
	return err
}
