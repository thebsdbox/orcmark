package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/orcmark/pkg/listener"
)

// This is the CLI for managing the Listening server, that will collect results from containers that have started.

// This is the listening port for incoming connections
var port int

// The path to a file to store results
var resultsPath string

// expected results to watch for
var expectedResults int

func init() {
	listenerCMD.Flags().IntVar(&port, "port", 8080, "The port to listen on")
	listenerCMD.Flags().StringVar(&resultsPath, "path", "", "The path to a file to store benchmark results")
	listenerCMD.Flags().IntVar(&expectedResults, "results", 0, "The number of results to watch for")

	// Add subcommands to the main application
	orcmarkCmd.AddCommand(listenerCMD)

}

var listenerCMD = &cobra.Command{
	Use:   "listener",
	Short: "Manage the listening server",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(loglevel))

		// Check that the port is > 1024 other wise root permissions are needed
		if port < 1024 {
			log.Fatalln("Chose a port that is higher than 1024")
		}
		listener.SetResultsPath(resultsPath)
		listener.SetResults(expectedResults)
		err := listener.StartListener(port)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
