package folding

import "fmt"

type plainFolder struct {
}

func (s plainFolder) Start(display string, _ bool) string {
	fmt.Println(display)
	return ""
}

func (s plainFolder) Stop(_ string) {
	return
}
