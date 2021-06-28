// Copyright (c) 2016, Daniel Martí <mvdan@mvdan.cc>
// github.com/mvdan/sh@v3.3.0/cmd/shfmt
// COPY of above (MIT) licensed project
// CHANGES:
// extracted file operations for a single file

package shfmt

import (
	"bytes"
	maybeio "github.com/google/renameio/maybe"
	"io"
	"mvdan.cc/sh/v3/fileutil"
	"mvdan.cc/sh/v3/syntax"
	"os"
)

func (s *shfmt) formatPath(path string, checkShebang bool) error {
	var readBuf bytes.Buffer
	var copyBuf = make([]byte, 32*1024)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	readBuf.Reset()
	if checkShebang {
		n, err := f.Read(copyBuf[:32])
		if err != nil {
			return err
		}
		if !fileutil.HasShebang(copyBuf[:n]) {
			return nil
		}
		readBuf.Write(copyBuf[:n])
	}
	if _, err := io.CopyBuffer(&readBuf, f, copyBuf); err != nil {
		return err
	}
	f.Close()
	return s.formatBytes(readBuf.Bytes(), path)
}

func (s *shfmt) formatBytes(src []byte, path string) error {
	var writeBuf bytes.Buffer

	prog, err := s.parser.Parse(bytes.NewReader(src), path)
	if err != nil {
		return err
	}
	writeBuf.Reset()
	s.printer.Print(&writeBuf, prog)
	res := writeBuf.Bytes()
	if !bytes.Equal(src, res) {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		perm := info.Mode().Perm()
		// TODO: support atomic writes on Windows?
		if err := maybeio.WriteFile(path, res, perm); err != nil {
			return err
		}
	}
	return nil
}

type shfmt struct {
	parser  *syntax.Parser
	printer *syntax.Printer
}
