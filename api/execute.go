package api

// TODO(tyler) this has repeated exec calls

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// ContinueOnError Make flag
var ContinueOnError = false

// Verbose to print report regardless of return code
var Verbose = false

// Quiet silence extra output
var Quiet = false

func runFiles(stdexec execCmd, filepath string, rep chan report) {
	// Format any file received as input.
	params := append(stdexec.frontparams, filepath)

	for _, par := range stdexec.endparam {
		params = append(params, par)
	}

	c := exec.Command(stdexec.bin, params...)

	var outB, errB bytes.Buffer

	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB

	outS := ""
	errS := ""

	err := c.Run()

	if err != nil || Verbose {
		outS = outB.String()
		errS = errB.String()
	}

	rep <- report{filepath, err, outS, errS}
}

func filter(list []string, expect bool, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range list {
		result := f(entry)
		if result == expect {
			filteredLists = append(filteredLists, entry)
		}
	}
	return filteredLists
}

func printReportErr(rep report) error {
	if rep.err == nil {
		fmt.Println(" ✔️  " + rep.file)
		if Verbose && !Quiet {
			fmt.Println("	" + rep.stdout)
		}
		return nil
	}

	fmt.Println(" ✖️  " + rep.file)

	if Quiet {
		return rep.err
	}

	fmt.Println("	std out : ")
	fmt.Println(rep.stdout)
	fmt.Println("	std err : ")
	fmt.Println(rep.stderr)

	if ContinueOnError {
		fmt.Println("	cmd err : ")
		fmt.Println("	" + rep.err.Error())
	} else {
		log.Fatal(rep.err)
	}

	return rep.err
}

func parallelLoop(ex execCmd, filteredFilepaths []string) error {
	report := make(chan report)

	if exec.Command("which", ex.bin).Run() != nil {
		mess := "Could not find " + ex.bin + ", possibly not installed"
		return errors.New(mess)
	}

	for _, filepath := range filteredFilepaths {
		go runFiles(ex, filepath, report)
	}

	var lasterr error

	for range filteredFilepaths {
		err := printReportErr(<-report)

		if err != nil {
			// Don't just return we still have channels to join
			lasterr = err
		}
	}

	return lasterr
}
