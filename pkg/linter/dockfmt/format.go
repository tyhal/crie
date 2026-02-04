package dockfmt

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/reteps/dockerfmt/lib"
)

func (l *LintDockFmt) format(filepath string) error {
	if l.IndentSize == 0 {
		l.IndentSize = defaultIndent
	}
	config := &lib.Config{
		IndentSize:      l.IndentSize,
		TrailingNewline: l.TrailingNewline,
		SpaceRedirects:  l.SpaceRedirects,
	}

	src, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read from file %s: %w", filepath, err)
	}

	lines := strings.SplitAfter(strings.TrimSuffix(string(src), "\n"), "\n")
	dst := lib.FormatFileLines(lines, config)

	if len(dst) == 0 {
		return errors.New("output empty invalid file")
	}

	err = os.WriteFile(filepath, []byte(dst), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filepath, err)
	}

	return nil
}
