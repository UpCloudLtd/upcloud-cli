package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

// StopCommand creates the "server stop" command
func StopCommand() commands.Command {
	return &stopCommand{
		BaseCommand: commands.New(
			"stop",
			"Stop a server",
			"upctl server stop 00cbe2f3-4cf9-408b-afee-bd340e13cdd8",
			"upctl server stop 00cbe2f3-4cf9-408b-afee-bd340e13cdd8 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0",
			"upctl server stop my_server",
			"upctl server stop --wait my_server",
		),
	}
}

type stopCommand struct {
	*commands.BaseCommand
	StopType string
	wait     config.OptionalBoolean
	resolver.CachingServer
	completion.Server
}

// InitCommand implements Command.InitCommand
func (s *stopCommand) InitCommand() {
	//XXX: findout what to do with risky params (timeout actions)
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.StopType, "type", defaultStopType, "The type of stop operation. Available: soft, hard")
	config.AddToggleFlag(flags, &s.wait, "wait", false, "Wait for server to be in stopped state before returning.")
	s.AddFlags(flags)
}

// Execute implements commands.MultipleArgumentCommand
func (s *stopCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("Stopping server %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.StopServer(&request.StopServerRequest{
		UUID:     uuid,
		StopType: s.StopType,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	if s.wait.Value() {
		waitForServerState(uuid, upcloud.ServerStateStopped, svc, logline, msg)
	} else {
		logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
		logline.MarkDone()
	}

	return output.OnlyMarshaled{Value: res}, nil
}
