package runner

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/linter"
)

func getName(lint linter.Linter) string {
	if lint == nil {
		return ""
	}
	return lint.Name()
}

// List to print all languages chkConf fmt and always commands
func (s *RunConfiguration) List() error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"language", "checker", "formatter", "associated files"})
	for langName, l := range s.Languages {

		err := table.Append([]string{langName, getName(l.Chk), getName(l.Fmt), l.Regex.String()})
		if err != nil {
			return err
		}
	}
	err := table.Render()
	return err
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

func (s *RunConfiguration) runLinter(cleanupGroup *sync.WaitGroup, name string, lintType LintType) (err error) {
	selectedLinter, err := s.Languages[name].GetLinter(lintType)
	if err != nil {
		return
	}
	toLog := log.WithFields(log.Fields{"lang": name, "type": lintType.String()})

	if selectedLinter == nil {
		skip := toLog.WithFields(log.Fields{"flag": "skip"})
		skip.Debug("there are no configurations associated for this action")
		return
	}

	// Get the match for this formatter's files.
	reg := s.Languages[name].Regex

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
func (s *RunConfiguration) Run(lintType LintType) (err error) {

	// Get an initial list of files to use
	s.loadFileList()

	errCount := 0

	currentLangs := s.Languages
	if s.Options.Only != "" {
		lang, err := s.GetLanguage(s.Options.Only)
		if err != nil {
			return err
		}
		currentLangs = map[string]*Language{s.Options.Only: lang}
	}

	var cleanupGroup sync.WaitGroup
	defer func() {
		// TODO bad linter implementations can cleanup forever with no timeout
		cleanupGroup.Wait()
	}()

	// Run every linter.
	for name := range currentLangs {
		err := s.runLinter(&cleanupGroup, name, lintType)
		if err != nil {
			log.Error(err)
			errCount++
			if !s.Options.Continue {
				break
			}
		}
	}

	if errCount > 0 {
		return errors.New(strconv.Itoa(errCount) + " language(s) failed while " + lintType.String() + "'ing \u26c8")
	}

	log.Println("\u26c5  " + lintType.String() + "'ing passed")
	return nil
}
