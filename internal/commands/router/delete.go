package router

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type deleteCommand struct {
	*commands.BaseCommand
	service service.Network
}

// DeleteCommand creates the "delete router" command
func DeleteCommand(service service.Network) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a router"),
		service:     service,
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getRouterArgCompletionFunction(s.service))
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return routerRequest{
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
		}.send(args)
	}
}
