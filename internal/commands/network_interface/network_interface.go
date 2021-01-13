package network_interface

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/ip_address"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

const maxNetworkInterfaceActions = 10

func NetworkInterfaceCommand() commands.Command {
	return &networkInterfaceCommand{commands.New("network", "Manage network interface")}
}

type networkInterfaceCommand struct {
	*commands.BaseCommand
}

func handleIpAddress(ipStrings []string) ([]request.CreateNetworkInterfaceIPAddress, error) {
	var ipAddresses []request.CreateNetworkInterfaceIPAddress
	for _, ipAddrStr := range ipStrings {
		t, err := ip_address.GetFamily(ipAddrStr)
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
