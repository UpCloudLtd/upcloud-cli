package networkinterface

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/ipaddress"
)

const maxNetworkInterfaceActions = 10

// BaseNetworkInterfaceCommand creates the "network-interface" command.
func BaseNetworkInterfaceCommand() commands.Command {
	return &networkInterfaceCommand{commands.New("network-interface", "Manage network interface")}
}

type networkInterfaceCommand struct {
	*commands.BaseCommand
}

func mapIPAddressesToRequest(ipStrings []string) ([]request.CreateNetworkInterfaceIPAddress, error) {
	if len(ipStrings) == 0 {
		return nil, nil
	}
	ipAddresses := make([]request.CreateNetworkInterfaceIPAddress, 0, len(ipStrings))
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
