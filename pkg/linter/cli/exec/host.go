package exec

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// HostExecutor runs CLI tools directly on the host operating system.
type HostExecutor struct {
}

// WillHost checks whether the given binary can be found on the host PATH.
func WillHost(bin string) error {
	_, err := exec.LookPath(bin)
	return err
}

// Setup performs any required initialization for host execution.
func (e *HostExecutor) Setup() error {
	return nil
}

// Exec runs the configured CLI tool on the host against the provided file.
func (e *HostExecutor) Exec(i Instance, filePath string, stdout io.Writer, stderr io.Writer) error {
	targetFilePath := filePath
	if i.ChDir {
		targetFilePath = filepath.Base(filePath)
	}

	params := append(i.Start, targetFilePath)
	params = append(params, i.End...)

	c := exec.Command(i.Bin, params...)
	if i.ChDir {
		c.Dir = filepath.Dir(filePath)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.Dir = cwd
	}
	c.Stdout = stdout
	c.Stderr = stderr

	return c.Run()
}

// Cleanup releases any resources allocated during host execution setup.
func (e *HostExecutor) Cleanup() error {
	return nil
}
