package linter

// TODO(tyler) this has repeated exec calls

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// ShowPass is a global that will turn on and off printing files that report success
var ShowPass = true

// Log simple takes all fields and pushes them to our using the default logger
func (rep *Report) Log() error {

	if rep.Err == nil {
		if ShowPass {
			fmt.Printf("\u2714 %v\n", rep.File)
		}
		return nil
	}

	fmt.Printf("\u274C %v\n", rep.File)
	log.WithFields(log.Fields{"type": "stdout"}).Info(rep.StdOut)
	log.WithFields(log.Fields{"type": "stderr"}).Error(rep.StdErr)

	return rep.Err
}
