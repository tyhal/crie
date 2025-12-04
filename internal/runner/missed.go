package runner

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// NoStandards runs all fmt exec commands in languages and in always fmt
func (s *RunConfiguration) NoStandards() error {

	// GetFiles files not used
	files, err := s.getFileList()
	if err != nil {
		return err
	}
	for _, standardizer := range s.Languages {
		files = Filter(files, false, standardizer.FileMatch.MatchString)
	}

	// GetFiles extensions or Filename(if no extension) and count occurrences
	dict := make(map[string]int)
	for _, str := range files {

		_, s := filepath.Split(str)

		for i := len(str) - 1; i >= 0 && !os.IsPathSeparator(str[i]); i-- {
			if str[i] == '.' {
				s = str[i:]
			}
		}

		dict[s] = dict[s] + 1
	}

	// Print dict in order
	output := map[int][]string{}
	var values []int
	for i, file := range dict {
		output[file] = append(output[file], i)
	}
	for i := range output {
		values = append(values, i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	// Print the top 10
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"extension", "count"})
	count := 10
	for _, i := range values {
		for _, s := range output[i] {
			err = table.Append([]string{s, strconv.Itoa(i)})
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
	err = table.Render()
	if err != nil {
		return err
	}
	return nil
}
