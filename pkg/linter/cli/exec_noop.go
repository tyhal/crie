package cli

import "io"

type noopExecutor struct{}

func (ne *noopExecutor) setup() error {
	return nil
}

func (ne *noopExecutor) exec(bin string, frontParams []string, filePath string, endParams []string, chdir bool, stdout io.Writer, stderr io.Writer) error {
	_, err := stdout.Write([]byte("stdout"))
	if err != nil {
		return err
	}
	_, err = stderr.Write([]byte("stderr"))
	if err != nil {
		return err
	}
	return nil
}

func (ne *noopExecutor) cleanup() error {
	return nil
}
