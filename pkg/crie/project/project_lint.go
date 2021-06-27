package project

import (
	log "github.com/sirupsen/logrus"
	run_flint "github.com/tyhal/crie/internal/flint"
)

var projDirs []string

// CheckProjects for missing files that should be included
func (s *LintConfiguration) CheckProjects() {
	if s.SingleLang == "" && s.IsRepo() {
		log.WithFields(log.Fields{"projects": len(projDirs)}).Info("required files")
		for _, dir := range projDirs {
			run_flint.RunFlint(dir)
		}
	} else {
		log.Info("not checking for any projects")
	}
}
