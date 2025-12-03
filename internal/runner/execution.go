// Package runner contains the logic to run the linters
package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/errchain"
	"github.com/tyhal/crie/pkg/folding"
	"github.com/tyhal/crie/pkg/linter"
)

func getName(lint linter.Linter) string {
	if lint == nil {
		return ""
	}
	return lint.Name()
}

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

//func (s *RunConfiguration) runLinter(cleanupGroup *sync.WaitGroup, name string, lintType LintType, list []string) (err error) {
//	selectedLinter, err := s.Languages[name].GetLinter(lintType)
//	if err != nil {
//		return
//	}
//	toLog := log.WithFields(log.Fields{"lang": name, "type": lintType.String()})
//
//	if selectedLinter == nil {
//		skip := toLog.WithFields(log.Fields{"flag": "skip"})
//		skip.Debug("there are no configurations associated for this action")
//		return
//	}
//
//	// GetFiles the match for this formatter's files.
//	reg := s.Languages[name].FileMatch
//
//	// find the associated files with our given regex to match.
//	associatedFiles := Filter(list, true, reg.MatchString)
//
//	// Skip language as no files found
//	if len(associatedFiles) == 0 {
//		return
//	}
//
//	cleanupGroup.Add(1)
//	defer func() { go selectedLinter.Cleanup(cleanupGroup) }()
//
//	err = selectedLinter.WillRun()
//	if err != nil {
//		toLog.Error(err)
//		return
//	}
//
//	toLog.WithFields(log.Fields{"files": len(associatedFiles)}).Infof("running %s", name)
//	reporter := linter.Runner{
//		ShowPass:      s.Options.Passes,
//		StrictLogging: s.Options.StrictLogging,
//	}
//	err = reporter.LintFileList(selectedLinter, associatedFiles)
//	return
//}

func (s *RunConfiguration) getRunningLanguages() (map[string]*Language, error) {
	currentLangs := s.Languages
	if s.Options.Only != "" {
		if lang, ok := s.Languages[s.Options.Only]; ok {
			currentLangs = map[string]*Language{s.Options.Only: lang}
		} else {
			return nil, fmt.Errorf("language %s not found", s.Options.Only)
		}
	}
	return currentLangs, nil
}

type LinterReady struct {
	linter linter.Linter
	lang   string
}

type FilematchReady struct {
	lang string
}

type WorkloadReady struct {
	linter linter.Linter
}

type Job struct {
	linter linter.Linter
	file   string
}

func (s *RunConfiguration) runLinters(lintType LintType, fileList []string) error {

	// every channel feeding up to workers matches worker count so that workers are priority
	maxWorkers := 16
	maxWorkloadBacklog := maxWorkers * 4

	currentLangs, err := s.getRunningLanguages()
	if err != nil {
		return err
	}

	lintFilesMatched := make(chan FilematchReady, maxWorkloadBacklog)
	lintReady := make(chan LinterReady, maxWorkloadBacklog)
	var linterFiles sync.Map

	var langStartupWG sync.WaitGroup
	var langFileMatchingWG sync.WaitGroup
	for langName, lang := range currentLangs {
		langFileMatchingWG.Go(func() {
			currLint, err := lang.GetLinter(lintType)
			if err != nil {
				log.Error(err)
				return
			}
			if currLint == nil {
				return
			}
			var hasMatched bool
			reg := lang.FileMatch
			var matched []string
			for _, file := range fileList {
				if reg.MatchString(file) {
					if !hasMatched {
						hasMatched = true
						langStartupWG.Go(func() {
							err := currLint.WillRun()
							if err != nil {
								log.Error(err)
								return
							}
							lintReady <- LinterReady{
								linter: currLint,
								lang:   langName,
							}
						})
					}
					matched = append(matched, file)
				}
			}
			if hasMatched {
				// TODO maybe a channel would be faster? this does batches
				linterFiles.Store(langName, matched)
				lintFilesMatched <- FilematchReady{
					lang: langName,
				}
			}
		})
	}

	var workersWG sync.WaitGroup
	jobs := make(chan Job, maxWorkloadBacklog)
	reports := make(chan linter.Report, maxWorkloadBacklog*16)
	for range maxWorkers {
		workersWG.Go(func() {
			for job := range jobs {
				if job.linter == nil {
					log.Error("oh no")
				} else {
					job.linter.Run(job.file, reports)
				}
			}
		})
	}

	// submit jobs
	go func() {
		active := make(map[string]WorkloadReady)
		submit := func(lang string, workload WorkloadReady) {
			if files, ok := linterFiles.Load(lang); ok {
				for _, file := range files.([]string) {
					jobs <- Job{linter: workload.linter, file: file}
				}
			}
		}
		for {
			select {
			case fm, ok := <-lintFilesMatched:
				if !ok {
					lintFilesMatched = nil
				}
				if workload, ok := active[fm.lang]; ok {
					submit(fm.lang, workload)
				} else {
					active[fm.lang] = WorkloadReady{}
				}
			case lr, ok := <-lintReady:
				if !ok {
					lintReady = nil
				}
				if workload, ok := active[lr.lang]; ok {
					workload.linter = lr.linter
					submit(lr.lang, workload)
				} else {
					active[lr.lang] = WorkloadReady{lr.linter}
				}
			}
			if lintFilesMatched == nil && lintReady == nil {
				break
			}
		}
		close(jobs)
	}()

	r := linter.Runner{
		ShowPass:      s.Options.Passes,
		StrictLogging: s.Options.StrictLogging,
	}
	var reportWG sync.WaitGroup
	reportWG.Go(func() {
		r.Folder = folding.New()
		for report := range reports {
			err := r.Log(&report)
			log.Error(err)
		}
	})

	langFileMatchingWG.Wait()
	langStartupWG.Wait()
	close(lintFilesMatched)
	close(lintReady)
	workersWG.Wait()
	close(reports)
	reportWG.Wait()

	return nil
}

// Run is the generic way to run everything based on the package configuration
func (s *RunConfiguration) Run(lintType LintType) error {
	fileList, err := s.getFileList()
	if err != nil {
		return errchain.From(err).Link("getting files")
	}
	err = s.runLinters(lintType, fileList)
	if err != nil {
		return err
	}
	log.Println("\u26c5  " + lintType.String() + "'ing passed")
	return nil
}
