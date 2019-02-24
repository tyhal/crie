package api

func getfmtexec(l language) execCmd {
	return l.fmt
}

// Fmt runs all fmt exec commands in languages and in always fmt
func Fmt() error {
	return stdrun("fmt", getfmtexec)
}
