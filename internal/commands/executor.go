package commands

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/ui"
	"os"
	"time"
)

// Executor represents the execution context for commands
type Executor interface {
	NewLogEntry(s string) *ui.LogEntry
	Update()
	WaitFor(waitFn func() error, timeout time.Duration) error
}

type executeResult struct {
	Job    int
	Result interface{}
	Error  error
}

type executorImpl struct {
	Config  *config.Config
	LiveLog *ui.LiveLog
}

func (e *executorImpl) WaitFor(waitFn func() error, timeout time.Duration) error {
	rv := make(chan error)
	timedOut := time.After(timeout)
	go func() {
		rv <- waitFn()
	}()
	select {
	case returned := <-rv:
		return returned
	case <-timedOut:
		return fmt.Errorf("timed out")
	}
}

func (e *executorImpl) NewLogEntry(message string) *ui.LogEntry {
	entry := ui.NewLogEntry(message)
	e.LiveLog.AddEntries(entry)
	return entry
}

// Update implements Executor
func (e *executorImpl) Update() {
	e.LiveLog.Render()
}

// NewExecutor creates the default Executor
func NewExecutor(c *config.Config) Executor {
	return &executorImpl{
		Config:  c,
		LiveLog: ui.NewLiveLog(os.Stderr, ui.LiveLogDefaultConfig),
	}
}
