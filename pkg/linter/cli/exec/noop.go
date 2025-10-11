package exec

import "io"

// NoopExecutor is a test/dummy executor that writes fixed strings to stdout and stderr.
type NoopExecutor struct{}

// Setup initializes the NoopExecutor (no-op).
func (ne *NoopExecutor) Setup() error {
	return nil
}

// Exec writes sample output to stdout and stderr and returns nil.
func (ne *NoopExecutor) Exec(_ Instance, _ string, stdout io.Writer, stderr io.Writer) error {
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

// Cleanup finalizes the NoopExecutor (no-op).
func (ne *NoopExecutor) Cleanup() error {
	return nil
}
