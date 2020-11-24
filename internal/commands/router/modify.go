package router

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	service service.Service
	req     request.ModifyRouterRequest
}

func ModifyCommand(service service.Service) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a router"),
		service:     service,
	}
}

func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.req.Name, "name", "", "Name of the router.")
	s.AddFlags(fs)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(router *upcloud.Router) interface{} {
				return &request.ModifyRouterRequest{UUID: router.UUID}
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ModifyRouterRequest).UUID },
				MaxActions:    maxRouterActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Modifying router",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyRouter(req.(*request.ModifyRouterRequest))
				},
			},
		}.Send(args)
	}
}
