package network

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/pflag"
)

const maxNetworkActions = 10

// BaseNetworkCommand creates the base "network" command.
func BaseNetworkCommand() commands.Command {
	return &networkCommand{
		BaseCommand: commands.New("network", "Manage network"),
	}
}

type networkCommand struct {
	*commands.BaseCommand
	resolver.CachingNetwork
}

// TODO: figure out a nicer way to do this..
func handleNetwork(in string) (*upcloud.IPNetwork, error) {
	result := &upcloud.IPNetwork{}
	var dhcp string
	var dhcpDefRout string
	var dns string

	args, err := commands.Parse(in)
	if err != nil {
		return nil, err
	}

	fs := &pflag.FlagSet{}
	fs.StringVar(&dns, "dhcp-dns", dns, "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")
	fs.StringVar(&result.Address, "address", result.Address, "Sets address space for the network.")
	fs.StringVar(&result.Family, "family", result.Address, "IP address family. Currently only IPv4 networks are supported.")
	fs.StringVar(&result.Gateway, "gateway", result.Gateway, "Gateway address given by the DHCP service. Defaults to first address of the network if not given.")
	fs.StringVar(&dhcp, "dhcp", dhcp, "Toggles DHCP service for the network.")
	fs.StringVar(&dhcpDefRout, "dhcp-default-route", dhcpDefRout, "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")

	err = fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if dhcp != "" {
		switch dhcp {
		case "true":
			result.DHCP = upcloud.FromBool(true)
		case "false":
			result.DHCP = upcloud.FromBool(false)
		default:
			return nil, fmt.Errorf("%s is an invalid value for dhcp, it can be true of false", dhcp)
		}
	}

	if dhcpDefRout != "" {
		if dhcpDefRout == "false" {
			result.DHCPDefaultRoute = upcloud.FromBool(false)
		}
		if dhcpDefRout == "true" {
			result.DHCPDefaultRoute = upcloud.FromBool(true)
		}
		return nil, fmt.Errorf("%s is an invalid value for dhcp default rout, it can be true of false", dhcp)
	}

	if dns != "" {
		result.DHCPDns = strings.Split(dns, ",")
	}

	return result, nil
}
