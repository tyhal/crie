package linter

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/folding"
)

type Reporter interface {
	Log(rep *Report) error
}

type StructuredReporter struct {
	ShowPass bool
	SrcOut   *log.Entry
	SrcErr   *log.Entry
	SrcInt   *log.Entry
}

func NewStructuredReporter(showPass bool) Reporter {
	return &StructuredReporter{
		ShowPass: showPass,
		SrcOut:   log.WithFields(log.Fields{"src": "stdout"}),
		SrcErr:   log.WithFields(log.Fields{"src": "stderr"}),
		SrcInt:   log.WithFields(log.Fields{"src": "internal"}),
	}
}

// Log simple takes all fields and pushes them to our using the default logger
func (r *StructuredReporter) Log(rep *Report) error {
	if rep.Err == nil {
		if r.ShowPass {
			log.Printf("pass %v", rep.File)
			r.SrcOut.Debug(rep.StdOut)
		}
	} else {
		log.Printf("fail %v", rep.File)
		var failedResultErr *FailedResultError
		if errors.As(rep.Err, &failedResultErr) {
			r.SrcErr.Error(rep.StdErr)
			r.SrcOut.Info(rep.StdOut)
			r.SrcInt.Debug(strings.NewReader(rep.Err.Error()), "toolerr", log.DebugLevel)
		} else {
			r.SrcInt.Error(strings.NewReader(rep.Err.Error()), "toolerr", log.ErrorLevel)
		}
	}
	return rep.Err
}

type StandardReporter struct {
	ShowPass bool
	SrcOut   *log.Entry
	SrcErr   *log.Entry
	SrcInt   *log.Entry
	Folder   folding.Folder
}

func NewStandardReporter(showPass bool) Reporter {
	return &StandardReporter{
		ShowPass: showPass,
		SrcOut:   log.WithFields(log.Fields{"src": "stdout"}),
		SrcErr:   log.WithFields(log.Fields{"src": "stderr"}),
		SrcInt:   log.WithFields(log.Fields{"src": "internal"}),
		Folder:   folding.New(),
	}
}

// Log simple takes all fields and pushes them to our using the default logger
func (r *StandardReporter) Log(rep *Report) error {
	if rep.Err == nil {
		if r.ShowPass {
			fmt.Printf("\u2714 %v\n", rep.File)
			r.SrcOut.Log(log.DebugLevel)
			fmt.Println(rep.StdOut)
		}
	} else {
		id, _ := r.Folder.Start(rep.File, "\u2716", false)
		var failedResultErr *FailedResultError
		if errors.As(rep.Err, &failedResultErr) {
			r.SrcErr.Log(log.ErrorLevel)
			fmt.Println(rep.StdErr)

			r.SrcOut.Log(log.InfoLevel)
			fmt.Println(rep.StdOut)

			r.SrcInt.Log(log.DebugLevel)
			fmt.Println(rep.Err)
		} else {
			r.SrcInt.Log(log.ErrorLevel)
			fmt.Println(rep.Err)
		}
		_ = r.Folder.Stop(id)
	}
	return rep.Err
}
