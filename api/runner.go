package api

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/api/linter"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// SingleLang if set will then tell languages to only use the given language
var SingleLang string

// GitDiff to use the git files instead of the entire tree
var GitDiff bool

func getLanguage(lang string) (linter.Language, error) {
	for _, standardizer := range languages {
		if standardizer.Name == lang {
			return standardizer, nil
		}
	}
	return languages[0], errors.New("language not found in configuration")
}

// NoStandards runs all fmt exec commands in languages and in always fmt
func NoStandards() {

	// Get files not used
	files := allFiles
	for _, standardizer := range languages {
		files = filter(files, false, standardizer.Match.MatchString)
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
	fmt.Println("Top Ten file types without standards applied to them")
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

// RunCrie is the generic way to run everything based on the packages configuration
func RunCrie() error {

	if GitDiff {
		allFiles = gitFiles
	}

	errCount := 0

	currentLangs := languages
	if SingleLang != "" {
		lang, err := getLanguage(SingleLang)
		if err != nil {
			return err
		}
		currentLangs = []linter.Language{lang}
	}

	// Run every linter.
	for _, l := range currentLangs {

		toLint := l.GetLinter(CurrentLinterType)
		toLog := log.WithFields(log.Fields{"lang": l.Name, "type": CurrentLinterType})

		if toLint == nil {
			toLog.Debug("there are no configurations associated for this action")
			continue
		}
		err := toLint.WillRun()
		if err != nil {
			toLog.Error(err.Error())
			errCount++
			continue
		}

		// Get the match for this formatter's files.
		reg := l.Match

		// filter the files to format based on given match and format them.
		filteredFilepaths := filter(allFiles, true, reg.MatchString)
		fmt.Println("❨ " + CurrentLinterType + " ❩ ➔ " + l.Name + " ❲" + strconv.Itoa(len(filteredFilepaths)) + "❳")

		err = parallelLoop(l, filteredFilepaths)

		if err != nil {
			toLog.Error(err.Error())
			errCount++
		}
	}

	if errCount > 0 {
		return errors.New("crie found " + strconv.Itoa(errCount) + " error(s) while " + CurrentLinterType + "'ing 🔍")
	}

	return nil
}
