package api

import (
	"fmt"
	"github.com/tyhal/crie/api/linter"
	"strings"
)

// CurrentLinterType is a global used to know if we are <chk|fmt|neither>
var CurrentLinterType = ""

var languages []linter.Language

// SetLinters is used to load in implemented linters from other packages
func SetLinters(l []linter.Language) {
	languages = l
}

// Version to print the current version of languages within crie
func Version() int {
	return len(languages)
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
	// TODO inteface that prints configs for each linter.Linter
}

// List to print all languages chkConf fmt and always commands
func List() {
	fmt.Println(" ~~~~~~~~~ ~~~~~~~~~")
	fmt.Println("\nLanguages:")
	for _, l := range languages {
		printLinter(l)
	}
	fmt.Println("\n~~~~~~~~~ ~~~~~~~~~")
}
