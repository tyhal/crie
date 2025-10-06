package crie

import (
	"errors"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/crie/linter"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)

// List to print all languages chkConf fmt and always commands
func (s *RunConfiguration) List() {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"language", "checker", "formatter", "associated files"})
	for langName, l := range s.Languages {
		table.Append([]string{langName, linter.GetName(l.Chk), linter.GetName(l.Fmt), l.Regex.String()})
	}
	table.Render()
}

// GetLanguage lets us query a language that might be in our projects' configuration
func (s *RunConfiguration) GetLanguage(lang string) (*linter.Language, error) {
	if language, exists := s.Languages[lang]; exists {
		return language, nil
	}
	return nil, errors.New("language not found in configuration")
}

// NoStandards runs all fmt exec commands in languages and in always fmt
func (s *RunConfiguration) NoStandards() {
	s.loadFileList()

	// Get files not used
	files := s.fileList
	for _, standardizer := range s.Languages {
		files = Filter(files, false, standardizer.Regex.MatchString)
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
	table.Header([]string{"extension", "count"})
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

func (s *RunConfiguration) runLinter(cleanupGroup *sync.WaitGroup, name string, lang *linter.Language) (err error) {
	selectedLinter, err := lang.GetLinter(s.lintType)
	if err != nil {
		return
	}
	toLog := log.WithFields(log.Fields{"lang": name, "type": s.lintType})

	if selectedLinter == nil {
		toLog.Debug("there are no configurations associated for this action")
		return
	}

	// Get the match for this formatter's files.
	reg := lang.Regex

	// filter the files to format based on given match and format them.
	filteredFilepaths := Filter(s.fileList, true, reg.MatchString)

	// Skip language as no files found
	if len(filteredFilepaths) == 0 {
		return
	}

	cleanupGroup.Add(1)
	defer func() { go selectedLinter.Cleanup(cleanupGroup) }()

	err = selectedLinter.WillRun()
	if err != nil {
		toLog.Error(err)
		return
	}

	toLog.WithFields(log.Fields{"files": len(filteredFilepaths)}).Infof("running %s", name)
	reporter := linter.Runner{
		ShowPass:      s.Options.Passes,
		StrictLogging: s.Options.StrictLogging,
	}
	err = reporter.LintFileList(selectedLinter, filteredFilepaths)
	return
}

// Run is the generic way to run everything based on the packages configuration
func (s *RunConfiguration) Run(lintType string) (err error) {
	s.lintType = lintType

	// Get an initial list of files to use
	s.loadFileList()

	errCount := 0

	currentLangs := s.Languages
	if s.Options.Only != "" {
		lang, err := s.GetLanguage(s.Options.Only)
		if err != nil {
			return err
		}
		currentLangs = map[string]*linter.Language{s.Options.Only: lang}
	}

	var cleanupGroup sync.WaitGroup
	defer func() {
		// TODO bad linter implementations can cleanup forever with no timeout
		cleanupGroup.Wait()
	}()

	// Run every linter.
	for name, lang := range currentLangs {
		err := s.runLinter(&cleanupGroup, name, lang)
		if err != nil {
			log.Error(err)
			errCount++
			if !s.Options.Continue {
				break
			}
		}
	}

	if errCount > 0 {
		return errors.New(strconv.Itoa(errCount) + " language(s) failed while " + s.lintType + "'ing \u26c8")
	}

	log.Println("\u26c5  " + s.lintType + "'ing passed")
	return nil
}
