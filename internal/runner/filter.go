package runner

// Filter to include (expect == true) or exclude (expect == false)
// the list of files (list) based on the filter function (f)
func Filter(list []string, expect bool, f func(string) bool) []string {
	filteredLists := make([]string, 0)
	for _, entry := range list {
		matched := f(entry)
		if matched == expect {
			filteredLists = append(filteredLists, entry)
		}
	}
	return filteredLists
}
