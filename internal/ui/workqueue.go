package ui

import (
	"io"
	"io/ioutil"
	"os"
	"time"
)

type WorkQueueConfig struct {
	NumTasks           int
	MaxConcurrentTasks int
	EnableUI           bool
}

func StartWorkQueue(cfg WorkQueueConfig, handler func(idx int, logEntry *LogEntry)) {
	var out io.Writer = os.Stdout
	if !cfg.EnableUI {
		out = ioutil.Discard
	}
	log := NewLiveLog(out, LiveLogDefaultConfig)

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
				startedTasks += 1
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
