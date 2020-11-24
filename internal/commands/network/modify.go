package network

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
	service service.Network
	req     request.ModifyNetworkRequest
}

func ModifyCommand(service service.Network) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify an ip address"),
		service:     service,
	}
}

func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.req.Name, "name", "", "Names the private network.")
	fs.StringVar(&s.req.Router, "router", "", "Change or clear the router attachment.")
	s.AddFlags(fs)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(network *upcloud.Network) interface{} {
				return &request.ModifyNetworkRequest{UUID: network.UUID}
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ModifyNetworkRequest).UUID },
				MaxActions:    maxNetworkActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Modifying network",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyNetwork(req.(*request.ModifyNetworkRequest))
				},
			},
		}.Send(args)
	}
}
