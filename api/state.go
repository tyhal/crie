package api

import (
	"fmt"
	"strings"
)

var standards = []language{
	lBash,
	lSh,
	lCpp,
	lCppheaders,
	lC,
	lCmake,
	lDocker,
	lGolang,
	lJavascript,
	lJSON,
	lPython,
	lMarkdown,
	lASCIIDoctor,
	lTerraform,
	lYML,
	lDockerCompose,
	lDoxygen,
}

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

// List to print all standards chk fmt and always commands
func List() {
	fmt.Println(" ~~~~~~~~~ ~~~~~~~~~")
	fmt.Println("\nLanguages:")
	for _, standard := range standards {
		fmt.Println("	" + standard.name)
		pprintCmd("❨ chk ❩ ", standard.chk.bin, standard.chk.frontparams, standard.chk.endparam)
		pprintCmd("❨ fmt ❩ ", standard.fmt.bin, standard.fmt.frontparams, standard.fmt.endparam)
	}
	fmt.Println("\n~~~~~~~~~ ~~~~~~~~~")
}
