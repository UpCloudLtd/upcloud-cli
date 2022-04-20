package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// StartCommand creates the "server start" command
func StartCommand() commands.Command {
	return &startCommand{
		BaseCommand: commands.New(
			"start",
			"Start a server",
			"upctl server start 00038afc-d526-4148-af0e-d2f1eeaded9b",
			"upctl server start 00038afc-d526-4148-af0e-d2f1eeaded9b 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0",
			"upctl server start my_server1",
		),
	}
}

type startCommand struct {
	*commands.BaseCommand
	completion.Server
	resolver.CachingServer
}

// InitCommand implements Command.InitCommand
func (s *startCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *startCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("Starting server %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.StartServer(&request.StartServerRequest{
		UUID: uuid,
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
