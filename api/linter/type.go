package linter

import "regexp"

type Linter interface {
	GetName() string
	GetReg() *regexp.Regexp
	Chk(file string, rep chan Report)
	Fmt(file string, rep chan Report)
}

type Report struct {
	File   string
	Err    error
	StdOut string
	StdErr string
}

type Lint func(file string, rep chan Report)
