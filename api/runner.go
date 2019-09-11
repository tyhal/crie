package api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// SingleLang if set will then tell standards to only use the given language
var SingleLang string

// GitDiff to use the git files instead of the entire tree
var GitDiff bool

func getLanguage(lang string) (language, error) {
	for _, standardizer := range standards {
		if standardizer.name == lang {
			return standardizer, nil
		}
	}
	return standards[0], errors.New("language not found in configuration")
}

// NoStandards runs all fmt exec commands in languages and in always fmt
func NoStandards() {

	// Get files not used
	files := allFiles
	for _, standardizer := range standards {
		files = filter(files, false, standardizer.match.MatchString)
	}

	// Get extensions or Filename(if no extension) and count occurrences
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
	fmt.Println("Top Ten file types without standards")
	count := 10
	for _, i := range values {
		for _, s := range output[i] {
			fmt.Printf("%s, %d\n", s, i)
			count--
			if count < 0 {
				return
			}
		}
	}
}

func stdrun(stdtype string, getexec execselec) error {

	if GitDiff {
		allFiles = gitFiles
	}

	errCount := 0

	runStandards := standards
	if SingleLang != "" {
		lang, err := getLanguage(SingleLang)
		if err != nil {
			return err
		}
		runStandards = []language{lang}
	}

	// Run every formatter.
	for _, standardizer := range runStandards {
		execC := getexec(standardizer)
		if execC.bin == "" {
			continue
		}

		// Get the match for this formatter's files.
		reg := standardizer.match

		// filter the files to format based on given match and format them.
		filteredFilepaths := filter(allFiles, true, reg.MatchString)
		fmt.Println("❨ " + stdtype + " ❩ ➔ " + standardizer.name + " ❲" + strconv.Itoa(len(filteredFilepaths)) + "❳")

		err := parallelLoop(execC, filteredFilepaths)

		if err != nil {
			fmt.Println(err.Error())
			errCount++
		}
	}

	if errCount > 0 {
		return errors.New("🔍 NGOT FAILED: " + strconv.Itoa(errCount) + " error(s) occurred while checking")
	}

	return nil
}
