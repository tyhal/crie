package crie

import (
	"io"
	"os"
)

// TODO remove helper package

// IsEmpty IsEmpty
func IsEmpty(name string) (bool, error) {
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

// RemoveIgnored Narrows down the list by returning only results that do not match the match in the config file
func RemoveIgnored(list []string, f func(string) bool) []string {
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

// Filter Filters
func Filter(list []string, expect bool, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range list {
		result := f(entry)
		if result == expect {
			filteredLists = append(filteredLists, entry)
		}
	}
	return filteredLists
}
