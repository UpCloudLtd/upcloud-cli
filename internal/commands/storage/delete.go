package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

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

func (s *deleteCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(storage *upcloud.Storage) (interface{}, error) {
				return &request.DeleteStorageRequest{UUID: storage.UUID}, nil
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
		}.Send(args)
	}
}
