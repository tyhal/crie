package imp

import (
	"bytes"
	"errors"
	"github.com/tyhal/crie/api/linter"
	"os"
	"os/exec"
)

type par []string

type execCmd struct {
	bin         string
	frontparams par
	endparam    par
}

func (e execCmd) WillRun() error {
	if exec.Command("which", e.bin).Run() != nil {
		return errors.New("could not find " + e.bin + ", possibly not installed")
	}
	return nil
}

func (e execCmd) Run(filepath string, rep chan linter.Report) {

	// Format any file received as input.
	params := append(e.frontparams, filepath)

	for _, par := range e.endparam {
		params = append(params, par)
	}

	c := exec.Command(e.bin, params...)

	var outB, errB bytes.Buffer

	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB

	outS := ""
	errS := ""

	err := c.Run()

	if err != nil {
		outS = outB.String()
		errS = errB.String()
	}

	rep <- linter.Report{filepath, err, outS, errS}
}
