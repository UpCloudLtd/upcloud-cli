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
func RestartCommand() commands.Command {
	return &restartCommand{
		BaseCommand: commands.New("restart", "Restart a server"),
	}
}

type restartCommand struct {
	*commands.BaseCommand
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
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))

	flags := &pflag.FlagSet{}

	flags.StringVar(&s.StopType, "stop-type", defaultStopType, "Restart type. Available: soft, hard")
	// TODO: reimplement? does not seem to make sense to automagically destroy servers if restart fails..
	// flags.StringVar(&s.params.TimeoutAction, "timeout-action", defaultRestartParams.TimeoutAction, "Action to take if timeout limit is exceeded. Available: destroy, ignore")
	flags.DurationVar(&s.Timeout, "timeout", defaultTimeout, "Server stop timeout in Go duration string\nExamples: 100ms, 1m10s, 3h")
	// TODO: reimplement? does not seem to be in use..
	// flags.IntVar(&s.params.Host, "host", defaultRestartParams.Host, "Use this to restart the VM on a specific host. Refers to value from host attribute. Only available for private cloud hosts")

	flags.BoolVar(&s.WaitForServerToStart, "wait", false, "Wait for server to start before exiting")

	s.AddFlags(flags)
}

func (s *restartCommand) MaximumExecutions() int {
	return maxServerActions
}

func (s *restartCommand) ArgumentMapper() (mapper.Argument, error) {
	return mapper.CachingServer(s.Config().Service.Server())
}

func (s *restartCommand) Execute(exec commands.Executor, uuid string) (output.Command, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("restarting server %v", uuid)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))
	res, err := svc.RestartServer(&request.RestartServerRequest{
		UUID:          uuid,
		StopType:      s.StopType,
		Timeout:       s.Timeout,
		TimeoutAction: "ignore",
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	if s.Config().GlobalFlags.Wait {
		logline.SetMessage(fmt.Sprintf("%s: waiting to restart", msg))
		if err := exec.WaitFor(serverStateWaiter(uuid, upcloud.ServerStateMaintenance, msg, svc, logline), s.Config().ClientTimeout()); err != nil {
			return nil, err
		}
		logline.SetMessage(fmt.Sprintf("%s: waiting to start", msg))
		if err := exec.WaitFor(serverStateWaiter(uuid, upcloud.ServerStateStarted, msg, svc, logline), s.Config().ClientTimeout()); err != nil {
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

func serverStateWaiter(uuid, state, msg string, service service.Server, logline *ui.LogEntry) func() error {
	return func() error {
		for {
			time.Sleep(100 * time.Millisecond)
			details, err := service.GetServerDetails(&request.GetServerDetailsRequest{UUID: uuid})
			if err != nil {
				return err
			}
			if details.State == state {
				return nil
			}
			logline.SetMessage(fmt.Sprintf("%s: waiting to start (%v)", msg, details.State))
		}
	}
}
