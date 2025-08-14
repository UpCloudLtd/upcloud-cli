package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/clierrors"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
	"github.com/gemalto/flume"
)

// Executor represents the execution context for commands
type Executor interface {
	Context() context.Context
	PushProgressUpdate(messages.Update)
	PushProgressStarted(msg string)
	PushProgressUpdateMessage(key, msg string)
	PushProgressSuccess(msg string)
	StopProgressLog()
	WaitFor(waitFn func() error, timeout time.Duration) error
	Server() service.Server
	Storage() service.Storage
	Network() service.Network
	Firewall() service.Firewall
	IPAddress() service.IPAddress
	Account() service.Account
	Token() service.Token
	All() internal.AllServices
	Debug(msg string, args ...interface{})
	WithLogger(args ...interface{}) Executor
	WithProgress(progress *progress.Progress) Executor
}

type executeResult struct {
	Job              int
	Result           output.Output
	Error            error
	ResolvedArgument resolvedArgument
}

type executorImpl struct {
	Config     *config.Config
	progress   *progress.Progress
	service    internal.AllServices
	logger     flume.Logger
	sigIntChan chan os.Signal
}

func (e executorImpl) WithLogger(args ...interface{}) Executor {
	e.logger = e.logger.With(args...)
	return &e
}

func (e executorImpl) WithProgress(progress *progress.Progress) Executor {
	e.progress = progress
	return &e
}

func (e *executorImpl) Context() context.Context {
	return e.Config.Context()
}

func (e *executorImpl) Debug(msg string, args ...interface{}) {
	e.logger.Debug(msg, args...)
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

func (e *executorImpl) PushProgressUpdate(update messages.Update) {
	err := e.progress.Push(update)
	if err != nil {
		e.Debug(fmt.Sprintf("Failed to push progress update: %s", err.Error()))
	}
}

func (e *executorImpl) PushProgressStarted(msg string) {
	e.PushProgressUpdate(messages.Update{
		Message: msg,
		Status:  messages.MessageStatusStarted,
	})
}

func (e *executorImpl) PushProgressUpdateMessage(key, msg string) {
	e.PushProgressUpdate(messages.Update{
		Key:     key,
		Message: msg,
	})
}

func (e *executorImpl) PushProgressSuccess(key string) {
	e.PushProgressUpdate(messages.Update{
		Key:    key,
		Status: messages.MessageStatusSuccess,
	})
}

func (e *executorImpl) StopProgressLog() {
	signal.Stop(e.sigIntChan) // prevent progress.Stop() from being called multiple times and panicking
	e.progress.Stop()
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

func (e executorImpl) IPAddress() service.IPAddress {
	return e.service
}

func (e executorImpl) Account() service.Account {
	return e.service
}

func (e executorImpl) Token() service.Token {
	return e.service
}

func (e executorImpl) All() internal.AllServices {
	return e.service
}

// NewExecutor creates the default Executor
func NewExecutor(cfg *config.Config, svc internal.AllServices, logger flume.Logger) Executor {
	executor := &executorImpl{
		Config:     cfg,
		progress:   progress.NewProgress(config.GetProgressOutputConfig()),
		logger:     logger,
		service:    svc,
		sigIntChan: make(chan os.Signal, 1),
	}
	executor.progress.Start()

	// Handle possible interrupts during execution
	signal.Notify(executor.sigIntChan, os.Interrupt)
	go func() {
		<-executor.sigIntChan

		// Cancel the app context
		cfg.Cancel()

		executor.PushProgressUpdate(messages.Update{
			Message: "Execution cancelled with interrupt signal",
			Details: "If you tried to create, modify, or delete resource(s), note that the operation might be already in progress. Please check the status of related resource(s).",
			Status:  messages.MessageStatusWarning,
		})
		executor.StopProgressLog()
		os.Exit(clierrors.InterruptSignalCode)
	}()

	return executor
}
