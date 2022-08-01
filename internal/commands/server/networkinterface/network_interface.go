package networkinterface

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

const maxNetworkInterfaceActions = 10

// BaseNetworkInterfaceCommand creates the "network-interface" command
func BaseNetworkInterfaceCommand() commands.Command {
	return &networkInterfaceCommand{commands.New("network-interface", "Manage network interface")}
}

type networkInterfaceCommand struct {
	*commands.BaseCommand
}

func (n *networkInterfaceCommand) BuildSubCommands(cfg *config.Config) {
	commands.BuildCommand(CreateCommand(), n.Cobra(), cfg)
	commands.BuildCommand(ModifyCommand(), n.Cobra(), cfg)
	commands.BuildCommand(DeleteCommand(), n.Cobra(), cfg)
}

func mapIPAddressesToRequest(ipStrings []string) ([]request.CreateNetworkInterfaceIPAddress, error) {
	var ipAddresses []request.CreateNetworkInterfaceIPAddress
	for _, ipAddrStr := range ipStrings {
		t, err := ipaddress.GetFamily(ipAddrStr)
		if err != nil {
			return nil, err
		}
		ip := request.CreateNetworkInterfaceIPAddress{
			Family:  t,
			Address: ipAddrStr,
		}
		ipAddresses = append(ipAddresses, ip)
	}
	return ipAddresses, nil
}
