package network

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

func DeleteCommand(service service.Network) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a network"),
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
				return &request.DeleteNetworkRequest{UUID: uuid}
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.DeleteNetworkRequest).UUID },
				MaxActions:    maxNetworkActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Deleting network",
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.service.DeleteNetwork(req.(*request.DeleteNetworkRequest))
				},
			},
		}.Send(args)
	}
}
