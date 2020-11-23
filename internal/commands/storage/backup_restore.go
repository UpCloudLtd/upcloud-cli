package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type restoreBackupCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  restoreBackupParams
}

type restoreBackupParams struct {
	request.RestoreBackupRequest
}

func RestoreBackupCommand(service service.Storage) commands.Command {
	return &restoreBackupCommand{
		BaseCommand: commands.New("restore", "Restore backup of a storage"),
		service:     service,
	}
}

func (s *restoreBackupCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(storage *upcloud.Storage) (interface{}, error) {
				req := s.params.RestoreBackupRequest
				req.UUID = storage.UUID
				return &req, nil
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestId:  func(in interface{}) string { return in.(*request.RestoreBackupRequest).UUID },
				MaxActions: maxStorageActions,
				ActionMsg:  "Restoring backup of storage",
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.service.RestoreBackup(req.(*request.RestoreBackupRequest))
				},
			},
		}.Send(args)
	}
}
