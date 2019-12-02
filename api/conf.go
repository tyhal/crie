package api

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// GlobalState to store relevant information
var GlobalState state

var projDirs []string

// allFiles the list of loaded files that need to be parsed
var allFiles []string

// gitFiles the list of loaded files that 'might' need to be parsed
var gitFiles []string

// CheckIgnores is to run against only the ignored files
var CheckIgnores = false

func newStdConf() {
	conf := conf{
		[]string{".git"},
		[]string{},
	}

	yamlOut, err := yaml.Marshal(conf)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(GlobalState.ConfName, yamlOut, 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New languages conf created: %s\nPlease view this and configure for your repo\n", GlobalState.ConfName)
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// Builds up a new list with allFiles matching the match in the config file
func getAllIgnored(allFiles []string, list []string, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range allFiles {
		result := f(entry)
		_, err := os.Stat(entry)
		if result && err == nil {
			filteredLists = append(filteredLists, entry)
		}
	}
	appendedList := append(list, filteredLists...)

	// Remove duplicates
	seen := make(map[string]struct{}, len(appendedList))
	j := 0
	for _, v := range appendedList {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		appendedList[j] = v
		j++
	}

	return appendedList[:j]
}

// Narrows down the list by returning only results that do not match the match in the config file
func removeIgnored(list []string, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range list {
		result := f(entry)
		_, err := os.Stat(entry)
		if !result && err == nil {
			filteredLists = append(filteredLists, entry)
		}
	}
	return filteredLists
}

func configureConf(files []string) []string {
	f, err := os.Open(GlobalState.ConfName)

	if err != nil {
		log.Fatal(err)
	}

	m := conf{}

	err = yaml.NewDecoder(f).Decode(&m)

	if err != nil {
		log.Fatal("Failed to parse (" + GlobalState.ConfName + "): " + err.Error())
	}

	var allFiles = files

	if CheckIgnores {
		files = nil // Clear the list - it'll have files added, rather than removed
	}

	for _, ignReg := range m.Ignore {
		reg, err := regexp.Compile(ignReg)

		if err != nil {
			log.Fatal(err)
		}

		if CheckIgnores {
			files = getAllIgnored(allFiles, files, reg.MatchString)
		} else {
			files = removeIgnored(files, reg.MatchString)
		}
	}

	// Add more project roots
	projDirs = m.ProjDirs
	projDirs = append(projDirs, ".")

	return files
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Initialise returns all valid files that have also been filtered by the config
func Initialise() {

	if Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Are we a repo?
	_, err := os.Stat(".git")
	GlobalState.IsRepo = err == nil

	// Work out where we are
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Create an initial file list
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			allFiles = append(allFiles, path)
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
	if _, err := os.Stat(GlobalState.ConfName); err == nil {
		allFiles = configureConf(allFiles)
	}

	if GlobalState.IsRepo {

		// If we are a repo without a configuration then force it upon the project
		if _, err := os.Stat(GlobalState.ConfName); err != nil {
			newStdConf()
			os.Exit(1)
		}

		var outB, errB bytes.Buffer

		c := exec.Command("git",
			par{"rev-list",
				"--no-merges",
				"--count",
				"HEAD"}...)

		c.Stdout = &outB

		if err = c.Run(); err != nil {
			log.Fatal(err)
		}

		// Produce string that will  query back all history or only 10 commits
		commitCntStr := strings.Split(outB.String(), "\n")[0]
		commitCnt, err := strconv.Atoi(commitCntStr)
		commitSlice := "HEAD~" + strconv.Itoa(min(commitCnt-1, 10)) + "..HEAD"

		args := par{"diff", "--name-only", commitSlice, "."}
		c = exec.Command("git", args...)

		c.Env = os.Environ()
		c.Stdout = &outB
		c.Stderr = &errB

		err = c.Run()

		if err != nil {
			fmt.Println("I noticed you are using git but I failed to get git diff")
			fmt.Println("... this is non-breaking (a-ok)")
			log.Debug(err.Error())
			log.WithFields(log.Fields{"type": "stdout"}).Debug(outB.String())
			log.WithFields(log.Fields{"type": "stderr"}).Debug(errB.String())
			gitFiles = allFiles
		} else {
			gitFiles = configureConf(strings.Split(outB.String(), "\n"))
		}
	}
}
