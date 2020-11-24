package storage

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
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
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		storages, err := s.service.GetStorages(&request.GetStoragesRequest{Access: upcloud.StorageAccessPrivate})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range storages.Storages {
			vals = append(vals, v.UUID, v.Title)
		}
		return commands.MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
	})
	s.SetPositionalArgHelp("<uuidOrTitle ...>")
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(storage *upcloud.Storage) (interface{}, error) {
				return &request.DeleteStorageRequest{UUID: storage.UUID}, nil
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
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
