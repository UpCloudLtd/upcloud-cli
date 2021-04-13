package network

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	networks []string
	name     string
	zone     string
	router   string
}

// make sure we implemnet the interface
var _ commands.Command = &createCommand{}

// CreateCommand creates the 'network create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a network", ""),
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
		"  dhcp-dns: array of strings \n"+
		"Usage: \n"+
		"	--ip-network 'address=94.23.112.143,\"dhcp-dns=<value1>,<value2>\",gateway=<gateway>,dhcp=true' \n"+
		"	--ip-network address=94.43.112.143/32,dhcp-dns=<value>\n")
	s.AddFlags(fs) // TODO(ana): replace usage with examples once the refactor is done.
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

	msg := fmt.Sprintf("creating network %v", s.name)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	req, err := s.buildRequest()
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	res, err := svc.CreateNetwork(req)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	return output.OnlyMarshaled{Value: res}, nil
}
