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

func (e *HostExecutor) Exec(i ExecInstance, filePath string, stdout io.Writer, stderr io.Writer) error {
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

func (e *HostExecutor) Cleanup() error {
	return nil
}
