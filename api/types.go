package api

type conf struct {
	Ignore   []string `yaml:"ignore"`
	ProjDirs []string `yaml:"proj_dirs"`
}

type state struct {
	IsRepo   bool
	ConfName string
}

type par []string
