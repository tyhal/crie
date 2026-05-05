package linter

import (
	"errors"
	"fmt"
	"strings"

	log "charm.land/log/v2"

	"github.com/tyhal/x/fold"
)

// Reporter is used to report results to the user
type Reporter interface {
	Log(rep *Report) error
}

// StructuredReporter is a Reporter that uses structured logging
type StructuredReporter struct {
	ShowPass bool
	SrcOut   *log.Logger
	SrcErr   *log.Logger
	SrcInt   *log.Logger
}

// NewStructuredReporter creates a new StructuredReporter
func NewStructuredReporter(showPass bool) Reporter {
	return &StructuredReporter{
		ShowPass: showPass,
		SrcOut:   log.With("src", "stdout"),
		SrcErr:   log.With("src", "stderr"),
		SrcInt:   log.With("src", "internal"),
	}
}

// Log simple takes all fields and pushes them to our using the default logFormat
func (r *StructuredReporter) Log(rep *Report) error {
	if rep.Err == nil {
		if r.ShowPass {
			log.With("target", rep.Target).Print("pass")
			r.SrcOut.Debug(rep.StdOut)
		}
	} else {
		if _, ok := errors.AsType[*FailedResultError](rep.Err); ok {
			// TODO do this better
			r.SrcErr.With("target", rep.Target).Error(rep.StdErr)
			r.SrcOut.With("target", rep.Target).Info(rep.StdOut)
			r.SrcInt.With("target", rep.Target).Debug(strings.NewReader(rep.Err.Error()), "toolerr", log.DebugLevel)
		} else {
			r.SrcInt.With("target", rep.Target).Error(strings.NewReader(rep.Err.Error()), "toolerr", log.ErrorLevel)
		}
	}

	return rep.Err
}

type logFormat struct {
	Entry *log.Logger
}

// Log prints the log message to stdout
func (l *logFormat) Log(level log.Level, args ...any) {
	if level >= log.GetLevel() {
		l.Entry.Log(level, "")
		fmt.Println(args...)
	}
}

// StandardReporter is a Reporter that uses simple logging
type StandardReporter struct {
	ShowPass bool
	SrcOut   logFormat
	SrcErr   logFormat
	SrcInt   logFormat
	fold.Folder
}

// NewStandardReporter creates a new StandardReporter
func NewStandardReporter(showPass bool) Reporter {
	return &StandardReporter{
		ShowPass: showPass,
		SrcOut:   logFormat{log.With("src", "stdout")},
		SrcErr:   logFormat{log.With("src", "stderr")},
		SrcInt:   logFormat{log.With("src", "internal")},
		Folder:   fold.New(),
	}
}

// Log simple takes all fields and pushes them to our using the default logFormat
func (r *StandardReporter) Log(rep *Report) error {
	if rep.Err == nil {
		if r.ShowPass {
			fmt.Printf("\u2714 %v\n", rep.Target)
			r.SrcOut.Log(log.DebugLevel, rep.StdOut)
		}
	} else {
		id, _ := r.Start(rep.Target, "\u2716", false)
		if _, ok := errors.AsType[*FailedResultError](rep.Err); ok {
			r.SrcErr.Log(log.ErrorLevel, rep.StdErr)
			r.SrcOut.Log(log.InfoLevel, rep.StdOut)
			r.SrcInt.Log(log.DebugLevel, rep.Err)
		} else {
			r.SrcInt.Log(log.ErrorLevel, rep.Err)
		}
		_ = r.Stop(id)
	}

	return rep.Err
}
