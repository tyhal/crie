package api

import (
	"fmt"
	"github.com/tyhal/crie/api/linter"
	"github.com/tyhal/crie/imp"
	"strings"
)

var CurrentLinterType = ""
var standards = imp.LanguageList

// Version to print the current version of standards within crie
func Version() int {
	return len(standards)
}

func pprintCmd(front string, bin string, frontparams []string, endparam []string) {
	if bin != "" {
		fmt.Println("		" + front + bin + " " + strings.Join(frontparams[:], " ") + " {file} " + strings.Join(endparam[:], " "))
	} else {
		fmt.Println("		⁉️  Not Implemented")
	}
}

func printLinter(l linter.Language) {
	fmt.Println("	" + l.Name)
	// TODO
	//pprintCmd("❨ chkConf ❩ ", l.chkConf.bin, l.chkConf.frontparams, l.chkConf.endparam)
	//pprintCmd("❨ fmtConf ❩ ", l.fmtConf.bin, l.fmtConf.frontparams, l.fmtConf.endparam)
}

// List to print all standards chkConf fmtConf and always commands
func List() {
	fmt.Println(" ~~~~~~~~~ ~~~~~~~~~")
	fmt.Println("\nLanguages:")
	for _, l := range standards {
		printLinter(l)
	}
	fmt.Println("\n~~~~~~~~~ ~~~~~~~~~")
}
