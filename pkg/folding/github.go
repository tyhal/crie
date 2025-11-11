package folding

import "fmt"

type githubFolder struct {
}

func (g githubFolder) Start(file string, msg string, _ bool) string {
	fmt.Printf("::error file=%s::%s %v\n", file, msg, file)
	fmt.Printf("::group::see logs\n")
	return file
}

func (g githubFolder) Stop(_ string) {
	fmt.Printf("::endgroup::\n")
}
