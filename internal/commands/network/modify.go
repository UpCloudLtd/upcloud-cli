package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	service  service.Network
	req      request.ModifyNetworkRequest
	networks []string
}

// ModifyCommand creates the "network modify" command
func ModifyCommand(service service.Network) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a network"),
		service:     service,
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getArgCompFn(s.service))
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.req.Name, "name", s.req.Name, "Names the private network.")
	fs.StringVar(&s.req.Router, "router", s.req.Router, "Change or clear the router attachment.")
	fs.StringArrayVar(&s.networks, "ip-network", s.networks, "The ip network with modified values. \n\n"+
		"Fields \n"+
		"  family: string \n"+
		"  gateway: string \n"+
		"  dhcp: true/false \n"+
		"  dhcp-default-route: true/false \n"+
		"  dhcp-dns: array of strings \n"+
		"Usage \n"+
		"	--ip-network dhcp-dns=<value1>,family=IPv4 \n"+
		" --ip-network 'dhcp=true,\"dhcp-dns=<value1>,<value2>\",family=IPv6'")
	s.AddFlags(fs) // TODO(ana): replace usage with examples once the refactor is done.
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		var networks []upcloud.IPNetwork
		for _, networkStr := range s.networks {
			network, err := handleNetwork(networkStr)
			if err != nil {
				return nil, err
			}
			if network.Family == "" {
				return nil, fmt.Errorf("family is required")
			}
			network.Address = ""
			networks = append(networks, *network)
		}
		s.req.IPNetworks = networks

		return networkRequest{
			BuildRequest: func(uuid string) interface{} {
				s.req.UUID = uuid
				return &s.req
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
		}.send(args)
	}
}
