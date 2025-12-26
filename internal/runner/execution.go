// Package runner contains the logic to run the linters
package runner

import (
	"context"
	"fmt"
	"runtime/trace"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/internal/runner/orchestrator"
	"github.com/tyhal/crie/pkg/errchain"
	"github.com/tyhal/crie/pkg/linter"
)

func (s *RunConfiguration) getRunningLanguages() (map[string]*Language, error) {
	currentLangs := s.Languages
	if s.Options.Only != "" {
		if lang, ok := s.Languages[s.Options.Only]; ok {
			currentLangs = map[string]*Language{s.Options.Only: lang}
		} else {
			return nil, fmt.Errorf("language %s not found", s.Options.Only)
		}
	}
	return currentLangs, nil
}

func (s *RunConfiguration) runLinters(ctx context.Context, lintType LintType, fileList []string) error {
	defer trace.StartRegion(ctx, "Crie Lint").End()

	currentLangs, err := s.getRunningLanguages()
	if err != nil {
		return err
	}

	var r linter.Reporter
	if s.Options.StrictLogging {
		r = linter.NewStructuredReporter(s.Options.Passes)
	} else {
		r = linter.NewStandardReporter(s.Options.Passes)
	}

	// NOTE an (obvious) assumption is made that formatters need file locking while linters do not
	locking := lintType == LintTypeFmt

	orch := orchestrator.New(fileList, r, locking)
	waitForCompletion := orch.Start(ctx)
	defer waitForCompletion()

	for _, lang := range currentLangs {
		l := lang.GetLinter(lintType)
		if l == nil {
			continue
		}
		orch.Dispatchers.Go(func() {
			orch.Dispatcher(ctx, l, lang.FileMatch)
		})
	}

	return nil
}

// Run is the generic way to run everything based on the package configuration
func (s *RunConfiguration) Run(ctx context.Context, lintType LintType) error {
	defer trace.StartRegion(ctx, "Crie Run").End()
	fileList, err := s.getFileList()
	if err != nil {
		return errchain.From(err).Link("getting files")
	}
	err = s.runLinters(ctx, lintType, fileList)
	if err != nil {
		return err
	}
	log.Println("\u26c5  " + lintType.String() + "'ing passed")
	return nil
}
