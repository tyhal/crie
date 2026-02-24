package runner

import (
	"fmt"
	"io"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// NoStandards runs all fmt exec commands in matchers and in always fmt
func (s *RunConfiguration) NoStandards(w io.Writer) error {
	// GetFiles files not used
	files, err := s.getFileList()
	if err != nil {
		return err
	}

	counts := s.noCoverageStats(files)

	return printCoverageStats(w, counts)
}

func cellStyle(row, col int) lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1)
}

func printCoverageStats(w io.Writer, counts map[string]int) error {
	fm := flatten(counts)
	slices.SortFunc(fm, fm.cmpV(true))

	// Print the top 10
	t := table.New()
	t.StyleFunc(cellStyle)
	t.Headers("ext", "#")
	for _, kv := range fm[:min(len(fm), 10)] {
		t.Row(kv.k, fmt.Sprintf("%d", kv.v))
	}

	_, err := fmt.Fprintln(w, t.Render())
	return err
}

func (s *RunConfiguration) noCoverageStats(files []string) map[string]int {
	for _, standardizer := range s.NamedMatches {
		files = Filter(files, false, standardizer.FileMatch.MatchString)
	}

	extCount := make(map[string]int)
	for _, path := range files {
		ext := filepath.Ext(path)
		if ext == "" {
			ext = filepath.Base(path)
		}
		extCount[ext]++
	}

	return extCount
}
