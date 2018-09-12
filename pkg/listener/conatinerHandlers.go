package listener

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// The HTTP handlers that the benchmark containers will use to report they're up

// Captures a container being up
func upPost(w http.ResponseWriter, r *http.Request) {
	counter++
	postTime := time.Since(stopwatch)
	log.Printf("UP Message [%d] from %s after %s", counter, r.RemoteAddr, postTime.String())
	// Only append the results time if were logging to a file
	if resultsFile != "" {
		// Convert result from nanoseconds to milliseconds, then convert to a string
		results = append(results, strconv.FormatInt((postTime.Nanoseconds()/1000000), 10))
	}
	if len(results) == expectedResults {
		writeResults()
	}
	io.WriteString(w, "")
}
