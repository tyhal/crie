package cli

import "io"

type executor interface {
	setup() error
	exec(bin string, frontParams []string, filePath string, endParams []string, chdir bool, stdout io.Writer, stderr io.Writer) error
	cleanup() error
}
