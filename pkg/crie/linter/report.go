package linter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

// ShowPass is a global that will turn on and off printing files that report success
var ShowPass = true

// StrictLogging allows to print straight to stdout for prettier formatting if it is off
var StrictLogging = false

func logConditional(reader io.Reader, typeField string, level log.Level) {
	field := log.WithFields(log.Fields{"type": typeField})
	if StrictLogging || level == log.DebugLevel {
		field.Log(level, reader)
	} else {
		field.Log(level)
		fmt.Println(reader)
	}
}

// Log simple takes all fields and pushes them to our using the default logger
func (rep *Report) Log() error {
	if rep.Err == nil {
		if ShowPass {
			if StrictLogging {
				log.Printf("pass %v", rep.File)
			} else {
				fmt.Printf("\u2714 %v\n", rep.File)
			}
			logConditional(rep.StdOut, "stdout", log.DebugLevel)
		}
	} else {
		if StrictLogging {
			log.Printf("fail %v", rep.File)

		} else {
			fmt.Printf("\u274C %v\n\n", rep.File)
		}
		logConditional(rep.StdErr, "stderr", log.ErrorLevel)
		logConditional(rep.StdOut, "stdout", log.InfoLevel)
		logConditional(rep.StdOut, "toolerr", log.DebugLevel)
	}

	return rep.Err
}
