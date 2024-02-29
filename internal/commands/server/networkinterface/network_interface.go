package networkinterface

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

const maxNetworkInterfaceActions = 10

// BaseNetworkInterfaceCommand creates the "network-interface" command
func BaseNetworkInterfaceCommand() commands.Command {
	return &networkInterfaceCommand{commands.New("network-interface", "Manage network interface")}
}

type networkInterfaceCommand struct {
	*commands.BaseCommand
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
