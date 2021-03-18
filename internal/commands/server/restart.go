package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/mapper"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"time"
)

// RestartCommand creates the "server restart" command
func RestartCommand(service service.Server) commands.Command {
	return &restartCommand{
		BaseCommand: commands.New("restart", "Restart a server"),
		service:     service,
	}
}

type restartCommand struct {
	*commands.BaseCommand
	service              service.Server
	WaitForServerToStart bool
	StopType             string
	TimeoutAction        string
	Timeout              time.Duration
}

const defaultStopType = "soft"
const defaultTimeout = time.Duration(120) * time.Second

// InitCommand implements Command.InitCommand
func (s *restartCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.service))
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.StopType, "stop-type", defaultStopType, "Restart type\nAvailable: soft, hard")
	// TODO: reimplement? does not seem to make sense to automagically destroy servers if restart fails..
	// flags.StringVar(&s.TimeoutAction, "timeout-action", defaultTimeoutAction, "Action to take if timeout limit is exceeded\nAvailable: destroy, ignore")
	flags.DurationVar(&s.Timeout, "timeout", defaultTimeout, "Stop timeout in Go duration string\nExamples: 100ms, 1m10s, 3h")
	// TODO: reimplement? does not seem to be in use..
	// flags.IntVar(&s.Host, "host", defaultHost, "Use this to restart the VM on a specific host. Refers to value from host attribute. Only available for private cloud hosts")
	flags.BoolVar(&s.WaitForServerToStart, "wait", false, "Wait for server to start before exiting")
	s.AddFlags(flags)
}

func (s *restartCommand) MaximumExecutions() int {
	return maxServerActions
}

func (s *restartCommand) ArgumentMapper() (mapper.Argument, error) {
	return mapper.CachingServer(s.service)
}

func (s *restartCommand) Execute(exec commands.Executor, uuid string) (output.Command, error) {
	msg := fmt.Sprintf("restarting server %v", uuid)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	res, err := s.service.RestartServer(&request.RestartServerRequest{
		UUID:          uuid,
		StopType:      s.StopType,
		Timeout:       s.Timeout,
		TimeoutAction: "ignore",
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		//		logline.SetMessage(fmt.Sprintf("failed (%v)", err))
		return nil, err
	}
	if s.WaitForServerToStart {
		logline.SetMessage(fmt.Sprintf("%s: waiting to restart", msg))
		if err := exec.WaitFor(serverStateWaiter(uuid, upcloud.ServerStateMaintenance, msg, s.service, logline), s.Config().ClientTimeout()); err != nil {
			return nil, err
		}
		logline.SetMessage(fmt.Sprintf("%s: waiting to start", msg))
		if err := exec.WaitFor(serverStateWaiter(uuid, upcloud.ServerStateStarted, msg, s.service, logline), s.Config().ClientTimeout()); err != nil {
			return nil, err
		}
		// TODO: this seems to not work as expected as the backend will report started->maintenance->started->maintenance..
		// should be fixed i guess?
		logline.SetMessage(fmt.Sprintf("%s: server restarted", msg))
	} else {
		logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	}
	logline.MarkDone()
	return output.Marshaled{Value: res}, nil
}

func (s *restartCommand) NewParent() commands.NewCommand {
	return s.Parent().(commands.NewCommand)
}

func serverStateWaiter(uuid, state, msg string, service service.Server, logline *ui.LogEntry) func() error {
	return func() error {
		for {
			time.Sleep(100 * time.Millisecond)
			details, err := service.GetServerDetails(&request.GetServerDetailsRequest{UUID: uuid})
			if err != nil {
				return err
			}
			if details.State == upcloud.ServerStateStarted {
				return nil
			}
			logline.SetMessage(fmt.Sprintf("%s: waiting to start (%v)", msg, details.State))
		}
	}
}
