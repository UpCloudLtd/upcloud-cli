package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/resolver"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

// DeleteCommand creates the "server delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a server"),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
	deleteStorages bool
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.BoolVar(&s.deleteStorages, "delete-storages", false, "Delete storages that are attached to the server.")
	s.AddFlags(flags)
}

func (s *deleteCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("deleting server %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	var err error
	if s.deleteStorages {
		logline.SetMessage(fmt.Sprintf("%s: deleting server and related storages", msg))
		err = svc.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
			UUID: uuid,
		})
	} else {
		logline.SetMessage(fmt.Sprintf("%s: deleting server", msg))
		err = svc.DeleteServer(&request.DeleteServerRequest{
			UUID: uuid,
		})
	}
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))

	return output.None{}, nil
}
