package cli

import (
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
)

func msgLast(fields []string) {
	sort.Slice(fields, func(i, j int) bool {
		if fields[i] == "msg" {
			return false
		}
		if fields[j] == "msg" {
			return true
		}
		return fields[i] < fields[j]
	})
}

func setLogging() {
	if projectConfig.Log.Trace {
		log.SetLevel(log.TraceLevel)
	}
	if projectConfig.Log.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if projectConfig.Log.Quiet {
		log.SetLevel(log.FatalLevel)
	}
	if projectConfig.Log.JSON {
		log.SetFormatter(&log.JSONFormatter{})
		projectConfig.Lint.StrictLogging = true
	} else {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&log.TextFormatter{
			SortingFunc:      msgLast,
			DisableQuote:     true,
			DisableTimestamp: true,
			DisableSorting:   false,
		})
	}
}
