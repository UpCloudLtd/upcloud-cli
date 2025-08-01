package network

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/pflag"
)

const maxNetworkActions = 10

var Types = []string{
	upcloud.NetworkTypePublic,
	upcloud.NetworkTypeUtility,
	upcloud.NetworkTypePrivate,
}

// BaseNetworkCommand creates the base "network" command
func BaseNetworkCommand() commands.Command {
	return &networkCommand{
		BaseCommand: commands.New("network", "Manage networks"),
	}
}

type networkCommand struct {
	*commands.BaseCommand
	resolver.CachingNetwork
}

// InitCommand implements Command.InitCommand
func (n *networkCommand) InitCommand() {
	n.Cobra().Aliases = []string{"net"}
}

// TODO: figure out a nicer way to do this..
func handleNetwork(in string) (*upcloud.IPNetwork, error) {
	result := &upcloud.IPNetwork{}
	var dhcp string
	var dhcpDefaultRoute string
	var dns string

	args, err := commands.Parse(in)
	if err != nil {
		return nil, err
	}

	fs := &pflag.FlagSet{}
	fs.StringVar(&dns, "dhcp-dns", dns, "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")
	fs.StringVar(&result.Address, "address", result.Address, "Sets address space for the network.")
	fs.StringVar(&result.Family, "family", result.Family, "IP address family. Currently only IPv4 networks are supported.")
	fs.StringVar(&result.Gateway, "gateway", result.Gateway, "Gateway address given by the DHCP service. Defaults to first address of the network if not given.")
	fs.StringVar(&dhcp, "dhcp", dhcp, "Toggles DHCP service for the network.")
	fs.StringVar(&dhcpDefaultRoute, "dhcp-default-route", dhcpDefaultRoute, "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")

	err = fs.Parse(args)
	if err != nil {
		return nil, err
	}

	commands.Must(fs.SetAnnotation("dhcp-dns", commands.FlagAnnotationFixedCompletions, []string{"true", "false"}))
	commands.Must(fs.SetAnnotation("address", commands.FlagAnnotationNoFileCompletions, nil))

	if dhcp != "" {
		val, err := commands.BoolFromString(dhcp)
		if err != nil {
			return nil, fmt.Errorf("could not parse dhcp value: %w", err)
		}
		result.DHCP = *val
	}

	if dhcpDefaultRoute != "" {
		val, err := commands.BoolFromString(dhcpDefaultRoute)
		if err != nil {
			return nil, fmt.Errorf("could not parse dhcp-default-route value: %w", err)
		}
		result.DHCPDefaultRoute = *val
	}

	if dns != "" {
		result.DHCPDns = strings.Split(dns, ",")
	}

	return result, nil
}
