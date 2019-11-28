package api

import (
	"fmt"
	"github.com/pengwynn/flint/flint"
	"os"
)

func flintRun(path string) {
	println("checking " + path)
	flags := flint.Flags{
		RunReadme:        true,
		RunContributing:  true,
		RunLicense:       true,
		RunBootstrap:     true,
		RunTestScript:    true,
		RunChangelog:     false,
		RunCodeOfConduct: false,
	}
	linter := flint.Linter{}
	summary, err := linter.Run(&flint.LocalProject{Path: path}, &flags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if summary != nil {
		summary.Print(os.Stderr, true)
		sev := summary.Severity()
		if sev > 0 {
			os.Exit(sev)
		}
	}
}

// Chk runs all Chk exec commands in languages and in always Chk
func Chk() error {
	if SingleLang == "" && GlobalState.IsRepo {
		for _, dir := range projDirs {
			flintRun(dir)
		}
	}
	CurrentLinterType = "chk"
	return RunDefaults()
}
