package api

import (
	"github.com/pengwynn/flint/flint"
	log "github.com/sirupsen/logrus"
	"os"
)

func flintRun(path string) {
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
		log.Fatal(err)
	}
	if summary != nil {
		sev := summary.Severity()
		if sev > 0 {
			log.WithFields(log.Fields{"path": path}).Error("project failed checks")
			summary.Print(os.Stderr, true)
			os.Exit(sev)
		}
	}
}

// Chk runs all Chk exec commands in languages and in always Chk
func (s *ProjectLintConfiguration) Chk() error {
	if s.SingleLang == "" && s.IsRepo {
		for _, dir := range projDirs {
			flintRun(dir)
		}
	}
	s.LintType = "chk"
	return s.Run()
}
