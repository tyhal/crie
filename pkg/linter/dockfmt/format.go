package dockfmt

import (
	"errors"
	"os"
	"strings"

	"github.com/reteps/dockerfmt/lib"
	"github.com/tyhal/crie/pkg/errchain"
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
		return errchain.From(err).LinkF("failed to read from file %s", filepath)
	}

	lines := strings.SplitAfter(strings.TrimSuffix(string(src), "\n"), "\n")
	dst := lib.FormatFileLines(lines, config)

	if len(dst) == 0 {
		return errors.New("output empty invalid file")
	}

	err = os.WriteFile(filepath, []byte(dst), 0644)
	if err != nil {
		return errchain.From(err).LinkF("failed to write to file %s", filepath)
	}

	return nil
}
