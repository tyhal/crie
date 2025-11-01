package folding

import "fmt"

type gitlabFolder struct {
}

func (g gitlabFolder) Start(id string) {
	fmt.Printf("start gitlab span %s\n", id)
}

func (g gitlabFolder) Stop() {
	//TODO implement me
	panic("implement me")
}

func (g gitlabFolder) Log() {
	//TODO implement me
	panic("implement me")
}
