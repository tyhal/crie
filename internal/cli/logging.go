package cli

import (
	log "charm.land/log/v2"
	"github.com/spf13/cobra"
)

func setLogging(cmd *cobra.Command) {
	log.SetOutput(cmd.OutOrStdout())
	if projectConfig.Log.Trace {
		log.SetLevel(log.DebugLevel)
	}
	if projectConfig.Log.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if projectConfig.Log.Quiet {
		log.SetLevel(log.FatalLevel)
	}
	if projectConfig.Log.JSON {
		log.SetFormatter(log.JSONFormatter)
		projectConfig.Lint.StrictLogging = true
	} else {
		log.SetReportTimestamp(false)
	}
}
