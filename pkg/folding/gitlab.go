package folding

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

type gitlabFolder struct {
	io.Writer
}

// NewGitlab uses the Gitlab CI log syntax
func NewGitlab(w io.Writer) Folder {
	return &gitlabFolder{w}
}

func (g gitlabFolder) Start(file, msg string, open bool) (string, error) {
	id := fmt.Sprintf("%08x\n", rand.Int31())
	collapsed := "true"
	if open {
		collapsed = "false"
	}
	_, err := fmt.Fprintf(g, "\033[0Ksection_start:%d:%s[collapsed=%s]\r\033[0K%s %v\n",
		time.Now().Unix(),
		id,
		collapsed,
		msg,
		file)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (g gitlabFolder) Stop(id string) error {
	_, err := fmt.Fprintf(g, "\033[0Ksection_end:%d:%s\r\033[0K", time.Now().Unix(), id)
	return err
}
