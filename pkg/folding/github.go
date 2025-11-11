package folding

import "fmt"

type githubFolder struct {
}

func (g githubFolder) Start(file string, msg string, _ bool) string {
	fmt.Printf("::error file=%s::%s\n", file, file)
	fmt.Printf("::group::%s see logs\n", msg)
	return file
}

func (g githubFolder) Stop(_ string) {
	fmt.Printf("::endgroup::\n")
}
