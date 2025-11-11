package folding

import "fmt"

type githubFolder struct {
}

func (g githubFolder) Start(display string, _ bool) string {
	fmt.Printf("::group::%s\n", display)
	return display
}

func (g githubFolder) Stop(_ string) {
	fmt.Printf("::endgroup::\n")
}
