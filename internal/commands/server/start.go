package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"

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
	exec.PushProgressStarted(msg)

	res, err := svc.StartServer(&request.StartServerRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
