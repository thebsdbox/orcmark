package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// The Used to set the logging level
var loglevel int

var orcmarkCmd = &cobra.Command{
	Use:   "orcmark",
	Short: "Orchestrator benchmark",
}

func init() {
	// Global flag across all subcommands
	orcmarkCmd.PersistentFlags().IntVar(&loglevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
}

// Execute - starts the command parsing process
func Execute() {
	log.SetLevel(log.Level(loglevel))

	if err := orcmarkCmd.Execute(); err != nil {
		orcmarkCmd.Help()
		fmt.Println(err)
		os.Exit(1)
	}
}
