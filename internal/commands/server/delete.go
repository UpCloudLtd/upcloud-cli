package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// DeleteCommand creates the "server delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a server",
			"upctl server delete 00cbe2f3-4cf9-408b-afee-bd340e13cdd8",
			"upctl server delete 00cbe2f3-4cf9-408b-afee-bd340e13cdd8 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0",
			"upctl server delete my_server",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.StoppedServer
	deleteStorages config.OptionalBoolean
	stop           config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &s.deleteStorages, "delete-storages", false, "Delete storages that are attached to the server.")
	config.AddToggleFlag(flags, &s.stop, "stop", false, "Stop the server before deleting it. Equivalent to running `upctl server stop --type hard --wait` before the delete command.")
	s.AddFlags(flags)
}

func Delete(exec commands.Executor, uuid, state string, deleteStorages, stopServer bool) (output.Output, error) {
	if stopServer && state != upcloud.ServerStateStopped {
		_, err := stop(exec, uuid, "hard", true)
		if err != nil {
			return nil, err
		}
	}

	svc := exec.Server()
	msg := fmt.Sprintf("Deleting server %v", uuid)
	exec.PushProgressStarted(msg)

	var err error
	if deleteStorages {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Deleting server %v and attached storages", uuid))
		err = svc.DeleteServerAndStorages(exec.Context(), &request.DeleteServerAndStoragesRequest{
			UUID: uuid,
		})
	} else {
		err = svc.DeleteServer(exec.Context(), &request.DeleteServerRequest{
			UUID: uuid,
		})
	}
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	server, _ := s.GetCached(uuid)

	return Delete(exec, uuid, server.State, s.deleteStorages.Value(), s.stop.Value())
}
