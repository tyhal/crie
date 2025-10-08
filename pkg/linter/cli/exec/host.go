package exec

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type HostExecutor struct {
}

func WillHost(bin string) error {
	_, err := exec.LookPath(bin)
	return err
}

func (e *HostExecutor) Setup() error {
	return nil
}

func (e *HostExecutor) Exec(bin string, frontParams []string, filePath string, endParams []string, chdir bool, stdout io.Writer, stderr io.Writer) error {
	targetFilePath := filePath
	if chdir {
		targetFilePath = filepath.Base(filePath)
	}

	params := append(frontParams, targetFilePath)
	params = append(params, endParams...)

	c := exec.Command(bin, params...)
	if chdir {
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

func (e *HostExecutor) Cleanup() error {
	return nil
}
