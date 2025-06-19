package crie

import (
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/flint/flint"
	"os"
)

// CheckProjectRequirements Ensures the project is Licensed and has the minimal documentation
func (s *RunConfiguration) CheckProjectRequirements(path string) {
	flags := flint.Flags{
		RunReadme:        true,
		RunContributing:  true,
		RunLicense:       true,
		RunBootstrap:     false,
		RunTestScript:    false,
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
