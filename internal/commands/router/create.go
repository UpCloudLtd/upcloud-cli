package router

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	service *service.Service
	params  createParams
}

func CreateCommand(service *service.Service) commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a router"),
		service:     service,
	}
}

type createParams struct {
	req     request.CreateRouterRequest
	routers []string
}

func (s *createCommand) InitCommand() {
	s.params.req = request.CreateRouterRequest{}
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.params.req.Name, "name", s.params.req.Name, "Name of the router.")
	s.AddFlags(fs)
}

func (s *createCommand) BuildRequest() (*request.CreateRouterRequest, error) {
	return &s.params.req, nil
}

func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		req, err := s.BuildRequest()
		if err != nil {
			return nil, err
		}
		var requests []interface{}
		requests = append(requests, req)

		return ui.HandleContext{
			RequestID:     func(in interface{}) string { return in.(*request.CreateRouterRequest).Name },
			ResultUUID:    getRouterUuid,
			MaxActions:    maxRouterActions,
			InteractiveUI: s.Config().InteractiveUI(),
			ActionMsg:     "Creating router",
			Action: func(req interface{}) (interface{}, error) {
				return s.service.CreateRouter(req.(*request.CreateRouterRequest))
			},
		}.HandleAction(requests)
	}
}
