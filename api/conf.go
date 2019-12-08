package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/api/linter"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var projDirs []string

func (s *ProjectLintConfiguration) initialiseRepo() {
	// If we are a repo without a configuration then force it upon the project
	if _, err := os.Stat(s.ConfPath); err != nil {
		createFileSettings(s.ConfPath)
	}

	var outB, errB bytes.Buffer

	c := exec.Command("git",
		par{"rev-list",
			"--no-merges",
			"--count",
			"HEAD"}...)

	c.Stdout = &outB

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	// Produce string that will  query back all history or only 10 commits
	commitCntStr := strings.Split(outB.String(), "\n")[0]
	commitCnt, err := strconv.Atoi(commitCntStr)
	commitSlice := "HEAD~" + strconv.Itoa(min(commitCnt-1, s.GitDiff)) + "..HEAD"

	args := par{"diff", "--name-only", commitSlice, "."}
	c = exec.Command("git", args...)

	c.Env = os.Environ()
	c.Stdout = &outB
	c.Stderr = &errB

	err = c.Run()

	if err != nil {
		log.WithFields(log.Fields{"type": "stdout"}).Debug(outB.String())
		log.WithFields(log.Fields{"type": "stderr"}).Debug(errB.String())
		log.Fatal(err.Error())
	} else {
		s.gitFiles = s.loadFileSettings(strings.Split(outB.String(), "\n"))
	}
}

func (s *ProjectLintConfiguration) initialiseFileList() {

	// Work out where we are
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Create an initial file list
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			s.allFiles = append(s.allFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	empty, err := isEmpty(".")
	if err != nil {
		log.Fatal(err)
	}

	if empty {
		return
	}

	// If there is a config then parse the files through it
	if _, err := os.Stat(s.ConfPath); err == nil {
		s.allFiles = s.loadFileSettings(s.allFiles)
	}
}

// loadFileList returns all valid files that have also been filtered by the config
func (s *ProjectLintConfiguration) loadFileList() {
	// Are we a repo?
	_, err := os.Stat(".git")
	s.IsRepo = err == nil

	if s.GitDiff > 0 {
		if s.IsRepo {
			s.initialiseRepo()
		} else {
			log.Fatal("this is not a git repo you are in")
		}
	} else {
		s.initialiseFileList()
	}
}

// SetLanguages is used to load in implemented linters from other packages
func (s *ProjectLintConfiguration) SetLanguages(l []linter.Language) {
	s.Languages = l
}

// List to print all languages chkConf fmt and always commands
func (s *ProjectLintConfiguration) List() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"language", "checker", "formatter", "associated files"})
	for _, l := range s.Languages {
		table.Append([]string{l.Name, getName(l.Chk), getName(l.Fmt), l.Match.String()})
	}
	table.Render()
}

// GetLanguage lets us query a language that might be in our projects configuration
func (s *ProjectLintConfiguration) GetLanguage(lang string) (linter.Language, error) {
	for _, standardizer := range s.Languages {
		if standardizer.Name == lang {
			return standardizer, nil
		}
	}
	return linter.Language{}, errors.New("language not found in configuration")
}

// NoStandards runs all fmt exec commands in languages and in always fmt
func (s *ProjectLintConfiguration) NoStandards() {
	s.loadFileList()

	// Get files not used
	files := s.allFiles
	for _, standardizer := range s.Languages {
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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"extension", "count"})
	count := 10
	for _, i := range values {
		for _, s := range output[i] {
			table.Append([]string{s, strconv.Itoa(i)})
			count--
			if count < 0 {
				table.Render()
				return
			}
		}
	}
	table.Render()
}

// Run is the generic way to run everything based on the packages configuration
func (s *ProjectLintConfiguration) Run() error {

	// Get initial list of files to use
	s.loadFileList()

	// Use git list if we are told to use atleast 1 or more git revisions
	if s.GitDiff > 0 {
		s.allFiles = s.gitFiles
	}

	errCount := 0

	currentLangs := s.Languages
	if s.SingleLang != "" {
		lang, err := s.GetLanguage(s.SingleLang)
		if err != nil {
			return err
		}
		currentLangs = []linter.Language{lang}
	}

	// Run every linter.
	for _, l := range currentLangs {

		selectedLinter := l.GetLinter(s.LintType)
		toLog := log.WithFields(log.Fields{"lang": l.Name, "type": s.LintType})

		if selectedLinter == nil {
			toLog.Debug("there are no configurations associated for this action")
			continue
		}

		err := selectedLinter.WillRun()
		if err != nil {
			toLog.Error(err.Error())
			errCount++
			continue
		}

		// Get the match for this formatter's files.
		reg := l.Match

		// filter the files to format based on given match and format them.
		filteredFilepaths := filter(s.allFiles, true, reg.MatchString)
		fmt.Println("❨ " + s.LintType + " ❩ ➔ " + l.Name + " ❲" + strconv.Itoa(len(filteredFilepaths)) + "❳")

		err = LintFileList(selectedLinter, filteredFilepaths)

		if err != nil {
			toLog.Error(err.Error())
			errCount++
			if !s.ContinueOnError {
				break
			}
		}
	}

	if errCount > 0 {
		return errors.New("crie found " + strconv.Itoa(errCount) + " error(s) while " + s.LintType + "'ing 🔍")
	}

	return nil
}
