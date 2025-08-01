package network

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
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
	fs.StringVar(&s.zone, "zone", s.zone, namedargs.ZoneDescription("network"))
	fs.StringVar(&s.router, "router", s.router, "Add this network to an existing router.")
	// TODO: handle multiline flag doc (try nested flags)
	fs.StringArrayVar(&s.networks, "ip-network", s.networks, "A network interface for the server, multiple can be declared.\n\n "+
		"Fields: \n"+
		"  address: string \n"+
		"  family: string \n"+
		"  gateway: string \n"+
		"  dhcp: true/false \n"+
		"  dhcp-default-route: true/false \n"+
		"  dhcp-dns: array of strings")

	s.AddFlags(fs)
	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("ip-network"))
	for _, flag := range []string{"name", "ip-network"} {
		commands.Must(s.Cobra().RegisterFlagCompletionFunc(flag, cobra.NoFileCompletions))
	}
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("router", namedargs.CompletionFunc(completion.Router{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
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
	svc := exec.Network()

	msg := fmt.Sprintf("Creating network %v", s.name)
	exec.PushProgressStarted(msg)

	req, err := s.buildRequest()
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	res, err := svc.CreateNetwork(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
