package clitool

import (
	"bytes"
	"errors"
	"github.com/tyhal/crie/api/linter"
	"os"
	"os/exec"
	"regexp"
)

func (e Language) GetName() string {
	return e.name
}

func (e Language) GetReg() *regexp.Regexp {
	return e.match
}

func (e execCmd) willRun() error {
	if exec.Command("which", e.bin).Run() != nil {
		mess := "Could not find " + e.bin + ", possibly not installed"
		return errors.New(mess)
	}
	return nil
}

func (e execCmd) run(filepath string) linter.Report {

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

	return linter.Report{filepath, err, outS, errS}
}


func (e Language) Chk(filepath string, rep chan linter.Report) {
	rep <- e.chkConf.run(filepath)
}

func (e Language) Fmt(filepath string, rep chan linter.Report) {
	rep <- e.fmtConf.run(filepath)
}
