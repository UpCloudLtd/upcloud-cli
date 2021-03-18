package commands

import "github.com/UpCloudLtd/cli/internal/config"

// Executor represents the execution context for commands
type Executor interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	WaitFor(func() error)
}

type executorImpl struct {
	Config *config.Config
}

// Log implements Executor
func (e *executorImpl) Log(args ...interface{}) {
	panic("implement me")
}

// Logf implements Executor
func (e *executorImpl) Logf(format string, args ...interface{}) {
	panic("implement me")
}

// WaitFor implements Executor
func (e *executorImpl) WaitFor(f func() error) {
	panic("implement me")
}

// NewExecutor creates the default Executor
func NewExecutor(c *config.Config) Executor {
	return &executorImpl{
		Config: c,
	}
}
