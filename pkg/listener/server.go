package listener

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Array to keep track of results
var results []string

// Path to where results should be stored
var resultsFile string

// The amount of expected results before storing
var expectedResults int

// StartListener - This will start the listening service
func StartListener(port int) error {
	log.Infof("Starting the listening server on port [%d]", port)
	log.Infoln("Press ctrl+c to stop the listener")

	if resultsFile != "" {
		log.Infof("Logging results to [%s]", resultsFile)
	}

	// Set the stop watch
	stopwatch = time.Now()

	// Stopwatch functionality
	http.HandleFunc("/stopwatch", resetStopwatch)

	// Container has started
	http.HandleFunc("/up", upPost)

	// Each Serve is part of a goroutine so we don't need to do anything special
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return err
	}
	return nil
}

func writeResults() error {
	// If the file doesn't exist, create it, or append to the file
	file, err := os.OpenFile(resultsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(results)

	defer writer.Flush()
	results = nil
	return nil
}

// SetResultsPath - determines where results should be stored
func SetResultsPath(path string) {
	resultsFile = path
}

// SetResults -
func SetResults(results int) {
	expectedResults = results
}
