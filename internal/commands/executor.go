package commands

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"os"
	"time"
)

// Executor represents the execution context for commands
type Executor interface {
	NewLogEntry(s string) *ui.LogEntry
	Update()
	WaitFor(waitFn func() error, timeout time.Duration) error
	Server() service.Server
	Storage() service.Storage
	Network() service.Network
	Firewall() service.Firewall
	IPAddress() service.IpAddress
	Account() service.Account
	Plan() service.Plans
	All() internal.AllServices
}

type executeResult struct {
	Job    int
	Result output.Output
	Error  error
}

type executorImpl struct {
	Config  *config.Config
	LiveLog *ui.LiveLog
	service internal.AllServices
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

func (e executorImpl) Server() service.Server {
	return e.service
}

func (e executorImpl) Storage() service.Storage {
	return e.service
}

func (e executorImpl) Network() service.Network {
	return e.service
}

func (e executorImpl) Firewall() service.Firewall {
	return e.service
}

func (e executorImpl) IPAddress() service.IpAddress {
	return e.service
}

func (e executorImpl) Account() service.Account {
	return e.service
}

func (e executorImpl) Plan() service.Plans {
	return e.service
}

func (e executorImpl) All() internal.AllServices {
	return e.service
}

// NewExecutor creates the default Executor
func NewExecutor(cfg *config.Config, svc internal.AllServices) Executor {
	return &executorImpl{
		Config:  cfg,
		LiveLog: ui.NewLiveLog(os.Stderr, ui.LiveLogDefaultConfig),
		service: svc,
	}
}
