package api


// Fmt runs all fmtConf exec commands in languages and in always fmtConf
func Fmt() error {
	return stdrun("fmt")
}
