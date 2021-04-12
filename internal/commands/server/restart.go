package server

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
	resolver.CachingServer
	completion.Server
	WaitForServerToStart bool
	StopType             string
	TimeoutAction        string
	Timeout              time.Duration
}

// InitCommand implements Command.InitCommand
func (s *restartCommand) InitCommand() {
	flags := &pflag.FlagSet{}

	// TODO: reimplement? does not seem to make sense to automagically destroy
	// servers if restart fails..
	flags.StringVar(&s.StopType, "stop-type", defaultStopType, "Restart type. Available: soft, hard")
	s.AddFlags(flags)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *restartCommand) MaximumExecutions() int {
	return maxServerActions
}

// Execute implements Command.Execute
func (s *restartCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("restarting server %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.RestartServer(&request.RestartServerRequest{
		UUID:          uuid,
		StopType:      s.StopType,
		Timeout:       defaultRestartTimeout,
		TimeoutAction: "ignore",
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	// TODO: reimplmement
	/*	if s.Config().GlobalFlags.Wait {
		// TODO: this seems to not work as expected as the backend will report
		// started->maintenance->started->maintenance..
		logline.SetMessage(fmt.Sprintf("%s: waiting to restart", msg))
		if err := exec.WaitFor(
			serverStateWaiter(uuid, upcloud.ServerStateMaintenance, msg, svc, logline),
			s.Config().ClientTimeout(),
		); err != nil {
			return nil, err
		}

		logline.SetMessage(fmt.Sprintf("%s: waiting to start", msg))
		if err := exec.WaitFor(
			serverStateWaiter(uuid, upcloud.ServerStateStarted, msg, svc, logline),
			s.Config().ClientTimeout(),
		); err != nil {
			return nil, err
		}

		logline.SetMessage(fmt.Sprintf("%s: server restarted", msg))
	} else {*/
	logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	//}

	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
