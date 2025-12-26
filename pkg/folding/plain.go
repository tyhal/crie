package folding

import (
	"fmt"
	"io"
)

type plainFolder struct {
	io.Writer
}

// NewPlain is a basic implementation of the Folder interface
func NewPlain(w io.Writer) Folder {
	return &plainFolder{w}
}

func (p plainFolder) Start(file, msg string, _ bool) (string, error) {
	_, err := fmt.Fprintf(p, "%s %v\n", msg, file)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (p plainFolder) Stop(_ string) error {
	return nil
}
