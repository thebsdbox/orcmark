package listener

import (
	"io"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Ths stopwatch is used to time the incoming requests
var stopwatch time.Time

// Counter track the amount of responses that have occured
var counter uint16

// This is used to reset or start the stopwatch
func resetStopwatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("Stopwatch start request from %s", r.RemoteAddr)
	stopwatch = time.Now()
	resetCounter()
	io.WriteString(w, "")
	if resultsFile != "" {
		writeResults()
	}
}

func resetCounter() {
	counter = 0
}
