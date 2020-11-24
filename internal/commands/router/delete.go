package router

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type deleteCommand struct {
	*commands.BaseCommand
	service *service.Service
}

func DeleteCommand(service *service.Service) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a router"),
		service:     service,
	}
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(router *upcloud.Router) interface{} {
				return &request.DeleteRouterRequest{UUID: router.UUID}
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.DeleteRouterRequest).UUID },
				MaxActions:    maxRouterActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Deleting router",
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.service.DeleteRouter(req.(*request.DeleteRouterRequest))
				},
			},
		}.Send(args)
	}
}
