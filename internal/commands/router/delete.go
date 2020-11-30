package router

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type deleteCommand struct {
	*commands.BaseCommand
	service Router
}

func DeleteCommand(service Router) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a router"),
		service:     service,
	}
}

func (s *deleteCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(uuid string) interface{} {
				return &request.DeleteRouterRequest{UUID: uuid}
			},
			Service: s.service,
			Handler: ui.HandleContext{
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
