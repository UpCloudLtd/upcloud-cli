package network

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"strings"
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
	fs.StringVar(&s.params.req.Name, "name", s.params.req.Name, "Names the network.")
	fs.StringVar(&s.params.req.Zone, "zone", s.params.req.Zone, "The zone in which the network is configured.")
	fs.StringVar(&s.params.req.Router, "router", s.params.req.Router, "Add this network to an existing router.")
	fs.StringArrayVar(&s.params.networks, "ip-network", s.params.networks, "A network interface for the server")
	s.AddFlags(fs)
}

func handleNetwork(in string) (*upcloud.IPNetwork, error) {
	network := upcloud.IPNetwork{}
	var dhcp bool
	var dhcpDefRout bool
	var dds string
	args, err := commands.Parse(in)
	if err != nil {
		return nil, err
	}

	fs := &pflag.FlagSet{}
	fs.StringVar(&dds, "dns", "", "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")
	fs.StringVar(&network.Family, "family", network.Family, "IP address family. Currently only IPv4 networks are supported.")
	fs.StringVar(&network.Address, "address", network.Address, "Sets address space for the network.")
	fs.StringVar(&network.Gateway, "gateway", network.Gateway, "Gateway address given by the DHCP service. Defaults to first address of the network if not given.")
	fs.BoolVar(&dhcp, "dhcp", dhcp, "Toggles DHCP service for the network.")
	fs.BoolVar(&dhcpDefRout, "dhcp-default-route", dhcpDefRout, "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")

	err = fs.Parse(args)
	if err != nil {
		return nil, err
	}

	network.DHCP = upcloud.FromBool(dhcp)
	network.DHCPDefaultRoute = upcloud.FromBool(dhcpDefRout)
	if dds != "" {
		network.DHCPDns = strings.Split(dds, ",")
	}

	return &network, nil
}

func (s *createCommand) BuildRequest() (*request.CreateNetworkRequest, error) {
	var networks []upcloud.IPNetwork
	for _, networkStr := range s.params.networks {
		network, err := handleNetwork(networkStr)
		if err != nil {
			return nil, err
		}
		networks = append(networks, *network)
	}
	s.params.req.IPNetworks = networks

	return &s.params.req, nil
}

func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		req, err := s.BuildRequest()
		if err != nil {
			return nil, err
		}
		var requests []interface{}
		requests = append(requests, req)

		return ui.HandleContext{
			RequestID:     func(in interface{}) string { return in.(*request.CreateNetworkRequest).Name },
			ResultUUID:    getNetworkUuid,
			MaxActions:    maxNetworkActions,
			InteractiveUI: s.Config().InteractiveUI(),
			ActionMsg:     "Creating network",
			Action: func(req interface{}) (interface{}, error) {
				return s.service.CreateNetwork(req.(*request.CreateNetworkRequest))
			},
		}.HandleAction(requests)
	}
}
