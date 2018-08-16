package listener

import (
	"io"
	"log"
	"net/http"
	"time"
)

// The HTTP handlers that the benchmark containers will use to report they're up

// Captures a container being up
func upPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("UP Message from %s after %v", r.RemoteAddr, time.Since(stopwatch))
	io.WriteString(w, "")
}
