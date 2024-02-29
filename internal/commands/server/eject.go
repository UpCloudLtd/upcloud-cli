package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type ejectCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
	params ejectParams
}

type ejectParams struct {
	request.EjectCDROMRequest
}

func (s *ejectCommand) InitCommand() {
}

// EjectCommand creates the "server eject" command
func EjectCommand() commands.Command {
	return &ejectCommand{
		BaseCommand: commands.New("eject", "Eject a CD-ROM from the server", "upctl server eject my_server"),
	}
}

// Execute implements commands.MultipleArgumentCommand
func (s *ejectCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()
	msg := fmt.Sprintf("Ejecting CD-ROM from %v", uuid)
	exec.PushProgressStarted(msg)

	req := s.params.EjectCDROMRequest
	req.ServerUUID = uuid

	res, err := svc.EjectCDROM(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
