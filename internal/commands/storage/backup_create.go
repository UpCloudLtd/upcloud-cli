package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type createBackupCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  createBackupParams
}

type createBackupParams struct {
	request.CreateBackupRequest
}

func CreateBackupCommand(service service.Storage) commands.Command {
	return &createBackupCommand{
		BaseCommand: commands.New("create", "Create backup of a storage"),
		service:     service,
	}
}

var DefaultCreateBackupParams = &createBackupParams{
	CreateBackupRequest: request.CreateBackupRequest{},
}

func (s *createBackupCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))
	s.params = createBackupParams{CreateBackupRequest: request.CreateBackupRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", DefaultCreateBackupParams.Title, "A short, informational description.\n[Required]")

	s.AddFlags(flagSet)
}

func (s *createBackupCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Title == "" {
			return nil, fmt.Errorf("title is required")
		}

		return Request{
			BuildRequest: func(storage *upcloud.Storage) (interface{}, error) {
				req := s.params.CreateBackupRequest
				req.UUID = storage.UUID
				return &req, nil
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.CreateBackupRequest).UUID },
				ResultUUID:    getStorageDetailsUuid,
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxStorageActions,
				ActionMsg:     "Creating backup of storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.CreateBackup(req.(*request.CreateBackupRequest))
				},
			},
		}.Send(args)
	}
}
