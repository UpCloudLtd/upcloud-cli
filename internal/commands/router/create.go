package router

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	service service.Network
	params  createParams
}

// CreateCommand creates the "router create" command
func CreateCommand(service service.Network) commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a router"),
		service:     service,
	}
}

type createParams struct {
	req request.CreateRouterRequest
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.params.req = request.CreateRouterRequest{}
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.params.req.Name, "name", s.params.req.Name, "Name of the router. [Required]")
	s.AddFlags(fs)
}

func (s *createCommand) buildRequest() (*request.CreateRouterRequest, error) {
	return &s.params.req, nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.req.Name == "" {
			return nil, fmt.Errorf("name is required")
		}

		req, err := s.buildRequest()
		if err != nil {
			return nil, err
		}

		return ui.HandleContext{
			RequestID:     func(in interface{}) string { return in.(*request.CreateRouterRequest).Name },
			ResultUUID:    getRouterUUID,
			MaxActions:    maxRouterActions,
			InteractiveUI: s.Config().InteractiveUI(),
			ActionMsg:     "Creating router",
			Action: func(req interface{}) (interface{}, error) {
				return s.service.CreateRouter(req.(*request.CreateRouterRequest))
			},
		}.Handle(commands.ToArray(req))
	}
}
