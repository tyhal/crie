package api

import (
	"github.com/pengwynn/flint/flint"
	log "github.com/sirupsen/logrus"
	"os"
)

func runFlint(path string) {
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

// CheckProjects for missing files that should be included
func (s *ProjectLintConfiguration) CheckProjects() {
	if s.SingleLang == "" && s.IsRepo() {
		log.WithFields(log.Fields{"projects": len(projDirs)}).Info("required files")
		for _, dir := range projDirs {
			runFlint(dir)
		}
	} else {
		log.Info("not checking for any projects")
	}
}
