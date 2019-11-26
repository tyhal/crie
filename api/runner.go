package api

import (
	"errors"
	"fmt"
	"github.com/tyhal/crie/api/clitool"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// SingleLang if set will then tell standards to only use the given language
var SingleLang string

// GitDiff to use the git files instead of the entire tree
var GitDiff bool

func getLanguage(lang string) (clitool.Language, error) {
	for _, standardizer := range standards {
		if standardizer.GetName() == lang {
			return standardizer, nil
		}
	}
	return standards[0], errors.New("language not found in configuration")
}

// NoStandards runs all fmtConf exec commands in languages and in always fmtConf
func NoStandards() {

	// Get files not used
	files := allFiles
	for _, standardizer := range standards {
		files = filter(files, false, standardizer.GetReg().MatchString)
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

func stdrun(stdtype string) error {

	if stdtype != "chk" {
		return errors.New(stdtype + " not implemented")
	}

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
		runStandards = []clitool.Language{lang}
	}

	// Run every formatter.
	for _, standardizer := range runStandards {

		// TODO check if cmd exists
		//if execC.bin == "" {
		//	continue
		//}

		// Get the match for this formatter's files.
		reg := standardizer.GetReg()

		// filter the files to format based on given match and format them.
		filteredFilepaths := filter(allFiles, true, reg.MatchString)
		fmt.Println("❨ " + stdtype + " ❩ ➔ " + standardizer.GetName() + " ❲" + strconv.Itoa(len(filteredFilepaths)) + "❳")

		err := parallelLoop(standardizer, filteredFilepaths)

		if err != nil {
			fmt.Println(err.Error())
			errCount++
		}
	}

	if errCount > 0 {
		return errors.New("🔍 CRIE FAILED: " + strconv.Itoa(errCount) + " error(s) occurred while checking")
	}

	return nil
}
