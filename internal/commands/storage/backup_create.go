package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
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

// CreateBackupCommand creates the "storage backup create" command
func CreateBackupCommand(service service.Storage) commands.Command {
	return &createBackupCommand{
		BaseCommand: commands.New("create", "Create backup of a storage"),
		service:     service,
	}
}

var defaultCreateBackupParams = &createBackupParams{
	CreateBackupRequest: request.CreateBackupRequest{},
}

// InitCommand implements Command.InitCommand
func (s *createBackupCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getStorageArgumentCompletionFunction(s.service))
	s.params = createBackupParams{CreateBackupRequest: request.CreateBackupRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", defaultCreateBackupParams.Title, "A short, informational description.")

	s.AddFlags(flagSet)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *createBackupCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Title == "" {
			return nil, fmt.Errorf("title is required")
		}

		return storageRequest{
			BuildRequest: func(uuid string) (interface{}, error) {
				req := s.params.CreateBackupRequest
				req.UUID = uuid
				return &req, nil
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.CreateBackupRequest).UUID },
				ResultUUID:    getStorageDetailsUUID,
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxStorageActions,
				ActionMsg:     "Creating backup of storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.CreateBackup(req.(*request.CreateBackupRequest))
				},
			},
		}.send(args)
	}
}
