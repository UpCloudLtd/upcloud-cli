package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
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

// TemplatizeCommand creates the "storage templatise" command
// TODO: figure out consistent naming, one way or the other.
func TemplatizeCommand(service service.Storage) commands.Command {
	return &templatizeCommand{
		BaseCommand: commands.New("templatise", "Templatise a storage"),
		service:     service,
	}
}

var defaultTemplatizeParams = &templatizeParams{
	TemplatizeStorageRequest: request.TemplatizeStorageRequest{},
}

// InitCommand implements Command.InitCommand
func (s *templatizeCommand) InitCommand() error {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getStorageArgumentCompletionFunction(s.service))
	s.params = templatizeParams{TemplatizeStorageRequest: request.TemplatizeStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", defaultTemplatizeParams.Title, "A short, informational description.")

	s.AddFlags(flagSet)

	return nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *templatizeCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Title == "" {
			return nil, fmt.Errorf("title is required")
		}

		return storageRequest{
			BuildRequest: func(uuid string) (interface{}, error) {
				req := s.params.TemplatizeStorageRequest
				req.UUID = uuid
				return &req, nil
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.TemplatizeStorageRequest).UUID },
				ResultUUID:    getStorageDetailsUUID,
				MaxActions:    maxStorageActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Templatising storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.TemplatizeStorage(req.(*request.TemplatizeStorageRequest))
				},
			},
		}.send(args)
	}
}
