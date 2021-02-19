package ui

import (
	"io"
	"io/ioutil"
	"os"
	"time"
)

// WorkQueueConfig represents the settings for a workqueue
type WorkQueueConfig struct {
	NumTasks           int
	MaxConcurrentTasks int
	EnableUI           bool
}

// StartWorkQueue starts a work queue that calls handler with idx specifying the current index in the work queue and
// logEntry that should be modified by the handler function to log entries.
func StartWorkQueue(cfg WorkQueueConfig, handler func(idx int, logEntry *LogEntry)) {
	var out io.Writer = os.Stdout
	if !cfg.EnableUI {
		out = ioutil.Discard
	}
	log := NewLiveLog(out, liveLogDefaultConfig)

	chDone := make(chan struct{}, cfg.NumTasks)
	doWork := func(idx int) {
		e := NewLogEntry("")
		log.AddEntries(e)
		handler(idx, e)
		e.MarkDone()
		chDone <- struct{}{}
	}
	startedTasks := 0
	completedTasks := 0
	for {
		select {
		case <-chDone:
			completedTasks++
		default:
			log.Render()
			switch {
			case completedTasks == cfg.NumTasks:
				log.Render()
				return
			case startedTasks-completedTasks < cfg.MaxConcurrentTasks && startedTasks < cfg.NumTasks:
				go doWork(startedTasks)
				startedTasks++
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
