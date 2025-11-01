package folding

import "fmt"

type plainFolder struct {
}

func (p plainFolder) Start(id string) {
	fmt.Printf("start plain span %s\n", id)
}

func (p plainFolder) Stop() {
	//TODO implement me
	panic("implement me")
}

func (p plainFolder) Log() {
	//TODO implement me
	panic("implement me")
}
