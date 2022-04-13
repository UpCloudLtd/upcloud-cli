package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

// RestartCommand creates the "server restart" command
func RestartCommand() commands.Command {
	return &restartCommand{
		BaseCommand: commands.New(
			"restart",
			"Restart a server",
			"upctl server restart 00038afc-d526-4148-af0e-d2f1eeaded9b",
			"upctl server restart 00038afc-d526-4148-af0e-d2f1eeaded9b --stop-type hard",
			"upctl server restart my_server1 my_server2",
		),
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
	flags.StringVar(&s.StopType, "stop-type", defaultStopType, "The type of stop operation. Available: soft, hard")
	s.AddFlags(flags)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *restartCommand) MaximumExecutions() int {
	return maxServerActions
}

// Execute implements commands.MultipleArgumentCommand
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

	logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
