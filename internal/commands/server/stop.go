package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

// StopCommand creates the "server stop" command
func StopCommand() commands.Command {
	return &stopCommand{
		BaseCommand: commands.New("stop", "Stop a server", ""),
	}
}

type stopCommand struct {
	*commands.BaseCommand
	StopType string
	resolver.CachingServer
	completion.Server
}

// InitCommand implements Command.InitCommand
func (s *stopCommand) InitCommand() {
	//XXX: findout what to do with risky params (timeout actions)
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.StopType, "type", defaultStopType, "The type of stop operation. Available: soft, hard")
	s.AddFlags(flags)
}

// Execute implements commands.MultipleArgumentCommand
func (s *stopCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("stopping server %v", uuid)
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

	logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
