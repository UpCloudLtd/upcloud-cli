package network

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"strings"
)

// ShowCommand creates the "network show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show network details", ""),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingNetwork
	completion.Network
}

func (s *showCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	network, err := s.CachingNetwork.GetCached(arg)
	if err != nil {
		return nil, err
	}
	commonSection := output.CombinedSection{Key: "", Title: "", Contents: output.Details{
		Sections: []output.DetailSection{
			{
				Title: "Common",
				Rows: []output.DetailRow{
					{Title: "UUID:", Key: "uuid", Value: network.UUID, Color: ui.DefaultUUUIDColours},
					{Title: "Name:", Key: "name", Value: network.Name},
					{Title: "Router:", Key: "router", Value: network.Router},
					{Title: "Type:", Key: "type", Value: network.Type},
					{Title: "Zone:", Key: "zone", Value: network.Zone},
				},
			},
		},
	},
	}

	networkRows := make([]output.TableRow, 0)
	if len(network.IPNetworks) > 0 {
		for _, nip := range network.IPNetworks {
			networkRows = append(networkRows, output.TableRow{
				nip.Address,
				nip.Family,
				nip.DHCP.Bool(),
				nip.DHCPDefaultRoute.Bool(),
				strings.Join(nip.DHCPDns, " "),
			})
		}
	}
	networkSection := output.CombinedSection{
		Key:   "ip_networks",
		Title: "IP Networks:",
		Contents: output.Table{
			Columns: []output.TableColumn{
				{Key: "address", Header: "Address", Color: ui.DefaultAddressColours},
				{Key: "family", Header: "Family"},
				{Key: "dhcp", Header: "DHCP", Format: output.BoolFormat},
				{Key: "dhcp_default_route", Header: "DHCP Def Router", Format: output.BoolFormat},
				{Key: "dhcp_dns", Header: "DHCP DNS"},
			},
			Rows: networkRows,
		}}

	serverRows := make([]output.TableRow, 0)
	for _, server := range network.Servers {
		fetched, err := exec.Server().GetServerDetails(&request.GetServerDetailsRequest{UUID: server.ServerUUID})
		if err != nil {
			return nil, fmt.Errorf("error getting server %v details: %w", server.ServerUUID, err)
		}
		serverRows = append(serverRows, output.TableRow{
			fetched.UUID,
			fetched.Title,
			fetched.Hostname,
			fetched.State,
		})
	}
	serverSection := output.CombinedSection{
		Key:   "servers",
		Title: "Servers:",
		Contents: output.Table{
			Columns: []output.TableColumn{
				{Header: "UUID", Key: "uuid", Color: ui.DefaultUUUIDColours},
				{Header: "Title", Key: "title"},
				{Header: "Hostname", Key: "hostname"},
				{Header: "State", Key: "state"},
			},
			Rows: serverRows,
		}}

	return output.Combined{
		commonSection,
		networkSection,
		serverSection,
	}, nil
}
