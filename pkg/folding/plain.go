package folding

import "fmt"

type plainFolder struct {
}

func (s plainFolder) Start(file, msg string, _ bool) string {
	fmt.Printf("%s %v\n\n", msg, file)
	return ""
}

func (s plainFolder) Stop(_ string) {
	return
}
