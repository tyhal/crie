package folding

type plainFolder struct {
}

func (s plainFolder) Start(_ string, _ bool) string {
	return ""
}

func (s plainFolder) Stop(_ string) {
	return
}
