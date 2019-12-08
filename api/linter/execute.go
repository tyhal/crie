package linter

// TODO(tyler) this has repeated exec calls

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Log simple takes all fields and pushes them to our using the default logger
func (rep *Report) Log() error {

	if rep.Err == nil {
		fmt.Printf(" ✔️  %v\n", rep.File)
		return nil
	}

	fmt.Printf(" ✖️  %v\n", rep.File)

	log.WithFields(log.Fields{"type": "stdout"}).Info(rep.StdOut)
	log.WithFields(log.Fields{"type": "stderr"}).Error(rep.StdErr)
	log.WithFields(log.Fields{"type": "err"}).Debug(rep.Err.Error())

	return rep.Err
}
