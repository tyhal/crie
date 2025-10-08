package exec

import (
	"io"
	"os/exec"
	"path/filepath"
)

type hostExecutor struct {
}

func willHost(bin string) error {
	_, err := exec.LookPath(bin)
	return err
}

func (e *hostExecutor) setup() error {
	return nil
}

func (e *hostExecutor) exec(bin string, frontParams []string, filePath string, endParams []string, chdir bool, stdout io.Writer, stderr io.Writer) error {
	finalFilePath := filePath
	if chdir {
		finalFilePath = filepath.Join(finalFilePath, filepath.Dir(filePath))
	}

	params := append(frontParams, finalFilePath)
	params = append(params, endParams...)

	c := exec.Command(bin, params...)
	c.Dir = finalFilePath
	c.Stdout = stdout
	c.Stderr = stderr

	return c.Run()
}

func (e *hostExecutor) cleanup() error {
	return nil
}
