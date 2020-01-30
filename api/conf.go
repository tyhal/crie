package api

import (
	"errors"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/api/linter"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var projDirs []string

// IsRepo checks for a .git folder
func (s *ProjectLintConfiguration) IsRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}

// loadFileList returns all valid files that have also been filtered by the config
func (s *ProjectLintConfiguration) loadFileList() {

	fileList := []string{}
	var err error

	if s.IsRepo() {
		// If we are a repo without a configuration then force it upon the project
		if _, err := os.Stat(s.ConfPath); err != nil {
			createFileSettings(s.ConfPath)
		}

		if s.GitDiff > 0 {
			// Get files changed in last s.GitDiff commits
			fileList, err = s.fileListRepoChanged()
		} else {
			// Get all files in git repo
			fileList, err = s.fileListRepoAll()
		}
	} else {

		// Check if the user asked for git diffs when not in a repo
		if s.GitDiff > 0 {
			log.Fatal("This is not a git repo you are in")
		}

		// Generic grab all the files
		fileList, err = s.fileListAll()
	}
	if err != nil {
		log.Fatal(err.Error())
	} else {
		s.fileList = fileList
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
	files := s.fileList
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

func (s *ProjectLintConfiguration) tryLint(l linter.Language) error {
	selectedLinter := l.GetLinter(s.lintType)
	toLog := log.WithFields(log.Fields{"lang": l.Name, "type": s.lintType})

	if selectedLinter == nil {
		toLog.Debug("there are no configurations associated for this action")
		return nil
	}

	// Get the match for this formatter's files.
	reg := l.Match

	// filter the files to format based on given match and format them.
	filteredFilepaths := filter(s.fileList, true, reg.MatchString)

	// Skip language as no files found
	if len(filteredFilepaths) == 0 {
		return nil
	}

	err := selectedLinter.WillRun()
	if err != nil {
		toLog.Error(err.Error())
		return err
	}

	log.WithFields(log.Fields{"files": len(filteredFilepaths)}).Info(l.Name)

	err = LintFileList(selectedLinter, filteredFilepaths)
	selectedLinter.DidRun()
	if err != nil {
		toLog.Error(err.Error())
		return err
	}

	return nil
}

// Run is the generic way to run everything based on the packages configuration
func (s *ProjectLintConfiguration) Run(lintType string) error {

	s.lintType = lintType

	// Get initial list of files to use
	s.loadFileList()

	// XXX Set separate a separate packages setting based on our configuration
	linter.ShowPass = s.ShowPasses

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
		err := s.tryLint(l)
		if err != nil {
			errCount++
			if !s.ContinueOnError {
				break
			}
		}
	}

	if errCount > 0 {
		return errors.New("found " + strconv.Itoa(errCount) + " language(s) failed while " + s.lintType + "'ing \u26c8")
	}

	log.Println(s.lintType + "'ing passed \u26c5")
	return nil
}
