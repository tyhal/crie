package folding

import "fmt"

type githubFolder struct {
}

func (g githubFolder) Start(file string, msg string, _ bool) string {
	fmt.Printf("::error file=%s::%s\n", file, msg)
	fmt.Printf("::group::see more\n")
	return file
}

func (g githubFolder) Stop(_ string) {
	fmt.Printf("::endgroup::\n")
}
