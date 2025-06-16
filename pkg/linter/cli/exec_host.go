package cli

import "os/exec"

func (e *Lint) willHost() error {
	_, err := exec.LookPath(e.Bin)
	return err
}
