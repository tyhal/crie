// Package runner contains the logic to run the linters
package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/trace"
	"strings"

	log "charm.land/log/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/tyhal/crie/internal/runner/orchestrator"
	"github.com/tyhal/crie/pkg/linter"
)

var tracer = otel.Tracer("github.com/tyhal/crie/internal/runner")

func (s *RunConfiguration) getRunningLanguages() (NamedMatches, error) {
	currentLangs := s.NamedMatches
	if s.Options.Only != "" {
		if lang, ok := s.NamedMatches[s.Options.Only]; ok {
			// if we have the specific language, drop all others
			currentLangs = NamedMatches{s.Options.Only: lang}
		} else {
			return nil, fmt.Errorf("language '%s' not found", s.Options.Only)
		}
	}
	return currentLangs, nil
}

func (s *RunConfiguration) runLinters(ctx context.Context, lintType LintType, fileList []string) (err error) {
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
	if s.Options.GitHubAnnotations {
		r = linter.NewGitHubReporter(r, os.Stdout, githubDir())
	}

	// NOTE an (obvious) assumption is made that formatters need file locking while linters do not
	locking := lintType == LintTypeFmt

	orch := orchestrator.New(fileList, r, locking, !s.Options.Continue)
	ctx, waitForCompletion := orch.Start(ctx)
	defer func() { err = errors.Join(err, waitForCompletion()) }()

	for _, lang := range currentLangs {
		l := lang.GetLinter(lintType)
		if l == nil {
			continue
		}
		if e := orch.CreateDispatcher(ctx, l, lang.FileMatch, lang.GroupBy); e != nil {
			err = e
			return
		}
	}

	return
}

// githubDir is the working directory relative to the GitHub workspace, which annotated paths are relative to
func githubDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	rel, err := filepath.Rel(os.Getenv("GITHUB_WORKSPACE"), wd)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return ""
	}
	return filepath.ToSlash(rel)
}

// Run is the generic way to run everything based on the package configuration
func (s *RunConfiguration) Run(ctx context.Context, lintType LintType) (err error) {
	defer trace.StartRegion(ctx, "Crie Run").End()
	ctx, span := tracer.Start(ctx, lintType.String())
	if s.Options.Only != "" {
		span.SetAttributes(attribute.String("crie.only", s.Options.Only))
	}
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()
	l := log.With("type", lintType)
	fileList, err := s.getFileList()
	span.SetAttributes(attribute.Int("crie.files", len(fileList)))
	if err != nil {
		return fmt.Errorf("getting files: %w", err)
	}
	l.With("files", len(fileList)).Printf("running")
	err = s.runLinters(ctx, lintType, fileList)
	if err != nil {
		return fmt.Errorf("failed %s: %w", lintType, err)
	}
	l.Printf("passed")
	return nil
}
