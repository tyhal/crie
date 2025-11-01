package folding

import "fmt"

type githubFolder struct {
}

func (g githubFolder) Start(id string) {
	fmt.Printf("start github span %s\n", id)
}

func (g githubFolder) Stop() {
	fmt.Printf("stop github span\n")
}

func (g githubFolder) Log() {
	//TODO implement me
	panic("implement me")
}
