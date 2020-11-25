package network_interface

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

const maxNetworkInterfaceActions = 10

func NetworkInterfaceCommand() commands.Command {
	return &networkInterfaceCommand{commands.New("network-interface", "Manage network interface")}
}

type networkInterfaceCommand struct {
	*commands.BaseCommand
}

func handleIpAddress(ipStrings []string) ([]request.CreateNetworkInterfaceIPAddress, error) {
	var ipAddresses []request.CreateNetworkInterfaceIPAddress
	for _, ipAddrStr := range ipStrings {
		ip := request.CreateNetworkInterfaceIPAddress{}
		args, err := commands.Parse(ipAddrStr)
		if err != nil {
			return nil, err
		}

		fs := &pflag.FlagSet{}
		fs.StringVar(&ip.Address, "address", "", "A valid IP address within the IP space of the network. If not given, the next free address is selected.")
		fs.StringVar(&ip.Family, "family", "", "IP address family. Currently only IPv4 networks are supported.")

		err = fs.Parse(args)
		if err != nil {
			return nil, err
		}

		ipAddresses = append(ipAddresses, ip)
	}
	return ipAddresses, nil
}
