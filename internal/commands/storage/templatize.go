package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type templatizeCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  templatizeParams
}

type templatizeParams struct {
	request.TemplatizeStorageRequest
}

func TemplatizeCommand(service service.Storage) commands.Command {
	return &templatizeCommand{
		BaseCommand: commands.New("templatize", "Templatize a storage"),
		service:     service,
	}
}

var DefaultTemplatizeParams = &templatizeParams{
	TemplatizeStorageRequest: request.TemplatizeStorageRequest{},
}

func (s *templatizeCommand) InitCommand() {
	s.params = templatizeParams{TemplatizeStorageRequest: request.TemplatizeStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", DefaultTemplatizeParams.Title, "A short, informational description.")

	s.AddFlags(flagSet)
}

func (s *templatizeCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(storage *upcloud.Storage) (interface{}, error) {
				req := s.params.TemplatizeStorageRequest
				req.UUID = storage.UUID
				return &req, nil
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestId:     func(in interface{}) string { return in.(*request.TemplatizeStorageRequest).UUID },
				ResultUuid: getStorageDetailsUuid,
				MaxActions:    maxStorageActions,
				InteractiveUi: s.Config().InteractiveUI(),
				ActionMsg:     "Templatizing storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.TemplatizeStorage(req.(*request.TemplatizeStorageRequest))
				},
			},
		}.Send(args)
	}
}
