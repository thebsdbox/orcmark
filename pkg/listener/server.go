package listener

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// StartListener - This will start the listening service
func StartListener(port int) error {
	log.Infof("Starting the listening server on port [%d]", port)
	log.Infoln("Press ctrl+c to stop the listener")

	// Set the stop watch
	stopwatch = time.Now()

	// Add handlers
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
