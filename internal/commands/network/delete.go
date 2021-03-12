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

// DeleteCommand creates the 'network delete' command
func DeleteCommand(service service.Network) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a network"),
		service:     service,
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() error {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getArgCompFn(s.service))

	return nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return networkRequest{
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
		}.send(args)
	}
}
