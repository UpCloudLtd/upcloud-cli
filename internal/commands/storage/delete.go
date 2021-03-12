package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

// DeleteCommand creates the "storage delete" command
func DeleteCommand(service service.Storage) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a storage"),
		service:     service,
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	service service.Storage
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() error {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getStorageArgumentCompletionFunction(s.service))

	return nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return storageRequest{
			BuildRequest: func(uuid string) (interface{}, error) {
				return &request.DeleteStorageRequest{UUID: uuid}, nil
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.DeleteStorageRequest).UUID },
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxStorageActions,
				ActionMsg:     "Deleting",
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.service.DeleteStorage(req.(*request.DeleteStorageRequest))
				},
			},
		}.send(args)
	}
}
