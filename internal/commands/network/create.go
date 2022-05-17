package network

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	networks []string
	name     string
	zone     string
	router   string
}

// make sure we implement the interface
var _ commands.Command = &createCommand{}

// CreateCommand creates the 'network create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a network",
			`upctl network create --name "My Network" --zone pl-waw1 --ip-network address=10.0.1.0/24`,
			"upctl network create --name my_net --zone pl-waw1 --ip-network address=10.0.2.0/24,dhcp=true",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", s.name, "Names the network.")
	fs.StringVar(&s.zone, "zone", s.zone, "The zone in which the network is configured.")
	fs.StringVar(&s.router, "router", s.router, "Add this network to an existing router.")
	//XXX: handle multiline flag doc (try nested flags)
	fs.StringArrayVar(&s.networks, "ip-network", s.networks, "A network interface for the server, multiple can be declared.\n\n "+
		"Fields: \n"+
		"  address: string \n"+
		"  family: string \n"+
		"  gateway: string \n"+
		"  dhcp: true/false \n"+
		"  dhcp-default-route: true/false \n"+
		"  dhcp-dns: array of strings")
	s.AddFlags(fs)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return maxNetworkActions
}

func (s *createCommand) buildRequest() (*request.CreateNetworkRequest, error) {
	var networks []upcloud.IPNetwork
	for _, networkStr := range s.networks {
		network, err := handleNetwork(networkStr)
		if err != nil {
			return nil, err
		}

		if network.Address == "" {
			return nil, fmt.Errorf("address is required for ip-network")
		}
		derivedFamily, err := ipaddress.GetFamily(network.Address)
		if err != nil {
			return nil, err
		}
		if network.Family != "" && network.Family != derivedFamily {
			return nil, fmt.Errorf("family %s is invalid for address %s", network.Family, network.Address)
		}
		network.Family = derivedFamily
		networks = append(networks, *network)
	}
	return &request.CreateNetworkRequest{
		Name:       s.name,
		Zone:       s.zone,
		Router:     s.router,
		IPNetworks: networks,
	}, nil
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// TODO: should we, for example, accept name as the first argument instead of as a flag?
	if s.name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if s.zone == "" {
		return nil, fmt.Errorf("zone is required")
	}
	if len(s.networks) == 0 {
		return nil, fmt.Errorf("at least one IP network is required")
	}
	svc := exec.Network()

	msg := fmt.Sprintf("Creating network %v", s.name)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	req, err := s.buildRequest()
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	res, err := svc.CreateNetwork(req)
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
