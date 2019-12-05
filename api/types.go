package api

type conf struct {
	Ignore   []string `hcl:"ignore"`
	ProjDirs []string `hcl:"proj_dirs"`
}

type state struct {
	IsRepo   bool
	ConfName string
}

type par []string
