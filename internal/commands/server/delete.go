package server

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
)

// DeleteCommand creates the "server delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a server"),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	deleteStorages bool
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))
	flags := &pflag.FlagSet{}
	flags.BoolVar(&s.deleteStorages, "delete-storages", false, "Delete storages that are attached to the server.")
	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		svc := s.Config().Service.(service.Server)

		var action = func(uuid interface{}) (interface{}, error) {
			var err error
			if s.deleteStorages {
				err = svc.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
					UUID: uuid.(string),
				})
			} else {
				err = svc.DeleteServer(&request.DeleteServerRequest{
					UUID: uuid.(string),
				})
			}
			return nil, err
		}

		return Request{
			BuildRequest: func(uuid string) interface{} { return uuid },
			Service:      svc,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(string) },
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				ActionMsg:     "Deleting",
				Action:        action,
			},
		}.Send(args)
	}
}
