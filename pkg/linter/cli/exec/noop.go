package exec

import "io"

type NoopExecutor struct{}

func (ne *NoopExecutor) Setup() error {
	return nil
}

func (ne *NoopExecutor) Exec(_ ExecInstance, _ string, stdout io.Writer, stderr io.Writer) error {
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

func (ne *NoopExecutor) Cleanup() error {
	return nil
}
