package network

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	networks []string
	router   string
	name     string
	completion.Network
	resolver.CachingNetwork
}

// ModifyCommand creates the "network modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a network"),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", "", "Set name of the private network.")
	fs.StringVar(&s.router, "router", "", "Change or clear the router attachment.")
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

// Execute implements Command.Execute
func (s *modifyCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, errors.New("need network to modify")
	}
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

	msg := fmt.Sprintf("modifying network %v", arg)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := exec.Network().ModifyNetwork(&request.ModifyNetworkRequest{
		UUID:       arg,
		Name:       s.name,
		Zone:       "", // TODO: should this be implemented?
		Router:     s.router,
		IPNetworks: networks,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	return output.Marshaled{Value: res}, nil
}
