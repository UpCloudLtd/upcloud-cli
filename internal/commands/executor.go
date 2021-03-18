package commands

import (
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/ui"
	"os"
)

// Executor represents the execution context for commands
type Executor interface {
	NewLogEntry(s string) *ui.LogEntry
	Update()
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
