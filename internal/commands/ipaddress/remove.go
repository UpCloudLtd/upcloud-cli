package ipaddress

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type removeCommand struct {
	*commands.BaseCommand
	service service.IpAddress
}

// RemoveCommand creates the 'ip-address remove' command
func RemoveCommand(service service.IpAddress) commands.Command {
	return &removeCommand{
		BaseCommand: commands.New("remove", "Removes an ip address"),
		service:     service,
	}
}

// InitCommand implements Command.MakeExecuteCommand
func (s *removeCommand) InitCommand() error {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getArgCompFn(s.service))

	return nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *removeCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return ipAddressRequest{
			BuildRequest: func(address string) interface{} {
				return &request.ReleaseIPAddressRequest{IPAddress: address}
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ReleaseIPAddressRequest).IPAddress },
				MaxActions:    maxIPAddressActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Removing IP Address",
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.service.ReleaseIPAddress(req.(*request.ReleaseIPAddressRequest))
				},
			},
		}.send(args)
	}
}
