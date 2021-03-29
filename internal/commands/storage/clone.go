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

type cloneCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  cloneParams
}

type cloneParams struct {
	request.CloneStorageRequest
}

// CloneCommand creates the "storage clone" command
func CloneCommand(service service.Storage) commands.Command {
	return &cloneCommand{
		BaseCommand: commands.New("clone", "Clone a storage"),
		service:     service,
	}
}

var defaultCloneParams = &cloneParams{
	CloneStorageRequest: request.CloneStorageRequest{
		Tier: upcloud.StorageTierHDD,
	},
}

// InitCommand implements Command.InitCommand
func (s *cloneCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getStorageArgumentCompletionFunction(s.service))
	s.params = cloneParams{CloneStorageRequest: request.CloneStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Tier, "tier", defaultCloneParams.Tier, "The storage tier to use.")
	flagSet.StringVar(&s.params.Title, "title", defaultCloneParams.Title, "A short, informational description.")
	flagSet.StringVar(&s.params.Zone, "zone", defaultCloneParams.Zone, "The zone in which the storage will be created, e.g. fi-hel1.")

	s.AddFlags(flagSet)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *cloneCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Zone == "" || s.params.Title == "" {
			return nil, fmt.Errorf("title and zone are required")
		}

		return storageRequest{
			BuildRequest: func(uuid string) (interface{}, error) {
				req := s.params.CloneStorageRequest
				req.UUID = uuid
				return &req, nil
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.CloneStorageRequest).UUID },
				ResultUUID:    getStorageDetailsUUID,
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxStorageActions,
				ActionMsg:     "Cloning storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.CloneStorage(req.(*request.CloneStorageRequest))
				},
			},
		}.send(args)
	}
}
