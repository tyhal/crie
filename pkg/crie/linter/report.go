package linter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// ShowPass is a global that will turn on and off printing files that report success
var ShowPass = true

// StrictLogging allows to print straight to stdout for prettier formatting if it is off
var StrictLogging = false

// Log simple takes all fields and pushes them to our using the default logger
func (rep *Report) Log() error {

	if rep.Err == nil {
		if ShowPass {
			if StrictLogging {
				log.Printf("pass %v", rep.File)
			} else {
				fmt.Printf("\u2714 %v\n", rep.File)
			}
		}
		return nil
	}

	if StrictLogging {
		log.Printf("fail %v", rep.File)
		log.WithFields(log.Fields{"type": "toolerr"}).Debug(rep.Err)
		log.WithFields(log.Fields{"type": "stdout"}).Info(rep.StdOut)
		if rep.StdErr != nil {
			log.WithFields(log.Fields{"type": "stderr"}).Error(rep.StdErr)
		}
	} else {
		fmt.Printf("\u274C %v\n\n", rep.File)
		log.WithFields(log.Fields{"type": "toolerr"}).Debug(rep.Err)
		log.WithFields(log.Fields{"type": "stdout"}).Info()
		fmt.Println(rep.StdOut)
		if rep.StdErr != nil {
			log.WithFields(log.Fields{"type": "stderr"}).Error()
			fmt.Println(rep.StdErr)
		}
	}

	return rep.Err
}
