package runner

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
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

func printCoverageStats(w io.Writer, counts map[string]int) error {
	// Print dict in order
	output := map[int][]string{}
	var values []int
	for i, file := range counts {
		output[file] = append(output[file], i)
	}
	for i := range output {
		values = append(values, i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	// Print the top 10
	table := tablewriter.NewWriter(w)
	defer table.Close()
	table.Header([]string{"extension", "count"})
	count := 10
	for _, i := range counts {
		for _, ext := range output[i] {
			err := table.Append([]string{ext, strconv.Itoa(i)})
			if err != nil {
				return err
			}
			count--
			if count < 0 {
				err = table.Render()
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	err := table.Render()
	if err != nil {
		return err
	}
	return nil
}

func (s *RunConfiguration) noCoverageStats(files []string) map[string]int {
	for _, standardizer := range s.NamedMatches {
		files = Filter(files, false, standardizer.FileMatch.MatchString)
	}

	// GetFiles extensions or Filename(if no extension) and count occurrences
	extCount := make(map[string]int)
	for _, str := range files {

		_, s := filepath.Split(str)

		for i := len(str) - 1; i >= 0 && !os.IsPathSeparator(str[i]); i-- {
			if str[i] == '.' {
				s = str[i:]
			}
		}

		extCount[s] = extCount[s] + 1
	}

	return extCount
}
