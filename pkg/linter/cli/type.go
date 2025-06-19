package cli

import "io"

type executor interface {
	setup() error
	exec(bin string, frontParams []string, filePath string, endParams []string, chdir bool, stdout io.Writer, stderr io.Writer) error
	cleanup() error
}

// Lint defines a predefined command to run against a file
type Lint struct {
	Bin            string
	FrontPar       Par
	EndPar         Par
	ContainerImage string
	ChDir          bool
	executor       executor
	cleanedUp      chan error
}

// Par represents cli parameters
type Par []string
