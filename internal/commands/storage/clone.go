package storage

import (
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

func CloneCommand(service service.Storage) commands.Command {
	return &cloneCommand{
		BaseCommand: commands.New("clone", "Clone a storage"),
		service:     service,
	}
}

var DefaultCloneParams = &cloneParams{
	CloneStorageRequest: request.CloneStorageRequest{
		Tier: "hdd",
	},
}

func (s *cloneCommand) InitCommand() {
	s.params = cloneParams{CloneStorageRequest: request.CloneStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Tier, "tier", DefaultCloneParams.Tier, "The storage tier to use.")
	flagSet.StringVar(&s.params.Title, "title", DefaultCloneParams.Title, "A short, informational description.")
	flagSet.StringVar(&s.params.Zone, "zone", DefaultCloneParams.Zone, "The zone in which the storage will be created, e.g. fi-hel1.")

	s.AddFlags(flagSet)
}

func (s *cloneCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(storage *upcloud.Storage) interface{} {
				req := s.params.CloneStorageRequest
				req.UUID = storage.UUID
				return &req
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestId:     func(in interface{}) string { return in.(*request.CloneStorageRequest).UUID },
				ResultUuid:    getStorageDetailsUuid,
				InteractiveUi: s.Config().InteractiveUI(),
				MaxActions:    maxStorageActions,
				ActionMsg:     "Cloning storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.CloneStorage(req.(*request.CloneStorageRequest))
				},
			},
		}.Send(args)
	}
}
