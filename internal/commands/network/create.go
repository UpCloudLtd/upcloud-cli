package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/ip_address"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	service service.Network
	params  createParams
}

func CreateCommand(service service.Network) commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a network"),
		service:     service,
	}
}

type createParams struct {
	req      request.CreateNetworkRequest
	networks []string
}

func (s *createCommand) InitCommand() {
	s.params.req = request.CreateNetworkRequest{}
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.params.req.Name, "name", s.params.req.Name, "Names the network.[Required]")
	fs.StringVar(&s.params.req.Zone, "zone", s.params.req.Zone, "The zone in which the network is configured.[Required]")
	fs.StringVar(&s.params.req.Router, "router", s.params.req.Router, "Add this network to an existing router.")
	fs.StringArrayVar(&s.params.networks, "ip-network", s.params.networks, "A network interface for the server, multiple can be declared.\n\n "+
		"Fields \n\n"+
		"  address: string \n\n"+
		"  family: string \n\n"+
		"  gateway: string \n\n"+
		"  dhcp: true/false \n\n"+
		"  dhcp-default-route: true/false \n\n"+
		"  dhcp-dns: array of strings \n\n"+
		"Usage \n\n"+
		"	--ip-network 'address=94.23.112.143,\"dhcp-dns=<value1>,<value2>\",gateway=<gateway>,dhcp=true' \n\n"+
		"	--ip-network address=94.43.112.143/32,dhcp-dns=<value>\n\n"+
		"[Required]")
	s.AddFlags(fs)
}

func (s *createCommand) BuildRequest() (*request.CreateNetworkRequest, error) {
	var networks []upcloud.IPNetwork
	for _, networkStr := range s.params.networks {
		network, err := handleNetwork(networkStr)
		if err != nil {
			return nil, err
		}

		if network.Address == "" {
			return nil, fmt.Errorf("address is required for ip-network")
		}
		derivedFamily, err := ip_address.GetFamily(network.Address)
		if err != nil {
			return nil, err
		}
		if network.Family != "" && network.Family != derivedFamily {
			return nil, fmt.Errorf("family %s is invalid for address %s", network.Family, network.Address)
		}
		network.Family = derivedFamily
		networks = append(networks, *network)
	}
	s.params.req.IPNetworks = networks

	return &s.params.req, nil
}

func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.req.Name == "" {
			return nil, fmt.Errorf("name is required")
		}
		if s.params.req.Zone == "" {
			return nil, fmt.Errorf("zone is required")
		}

		req, err := s.BuildRequest()
		if err != nil {
			return nil, err
		}

		return ui.HandleContext{
			RequestID:     func(in interface{}) string { return in.(*request.CreateNetworkRequest).Name },
			ResultUUID:    getNetworkUuid,
			MaxActions:    maxNetworkActions,
			InteractiveUI: s.Config().InteractiveUI(),
			ActionMsg:     "Creating network",
			Action: func(req interface{}) (interface{}, error) {
				return s.service.CreateNetwork(req.(*request.CreateNetworkRequest))
			},
		}.Handle(commands.ToArray(req))
	}
}
