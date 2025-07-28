package network

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowCommand creates the "network show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show network details",
			"upctl network show 037a530b-533e-4cef-b6ad-6af8094bb2bc",
			"upctl network show 037a530b-533e-4cef-b6ad-6af8094bb2bc 0311480d-d0c0-4951-ab41-bf12097f5d3c",
			`upctl network show "My Network"`,
		),
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
	network, err := s.GetCached(arg)
	if err != nil {
		return nil, err
	}

	combined := output.Combined{
		output.CombinedSection{
			Key:   "",
			Title: "",
			Contents: output.Details{
				Sections: []output.DetailSection{
					{
						Title: "Common",
						Rows: []output.DetailRow{
							{Title: "UUID:", Key: "uuid", Value: network.UUID, Colour: ui.DefaultUUUIDColours},
							{Title: "Name:", Key: "name", Value: network.Name},
							{Title: "Router:", Key: "router", Value: network.Router},
							{Title: "Type:", Key: "type", Value: network.Type},
							{Title: "Zone:", Key: "zone", Value: network.Zone},
						},
					},
				},
			},
		},
		labels.GetLabelsSection(network.Labels),
	}

	if len(network.IPNetworks) > 0 {
		networkRows := make([]output.TableRow, 0)
		for _, nip := range network.IPNetworks {
			networkRows = append(networkRows, output.TableRow{
				nip.Address,
				nip.Family,
				nip.DHCP.Bool(),
				nip.DHCPDefaultRoute.Bool(),
				nip.DHCPDns,
			})
		}

		combined = append(combined, output.CombinedSection{
			Key:   "ip_networks",
			Title: "IP Networks:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "address", Header: "Address", Colour: ui.DefaultAddressColours},
					{Key: "family", Header: "Family"},
					{Key: "dhcp", Header: "DHCP", Format: format.Boolean},
					{Key: "dhcp_default_route", Header: "DHCP Def Router", Format: format.Boolean},
					{Key: "dhcp_dns", Header: "DHCP DNS", Format: formatShowDHCPDNS},
				},
				Rows: networkRows,
			},
		})
	}

	if len(network.Servers) > 0 {
		serverRows := make([]output.TableRow, 0)
		for _, server := range network.Servers {
			fetched, err := exec.Server().GetServerDetails(exec.Context(), &request.GetServerDetailsRequest{UUID: server.ServerUUID})
			if err != nil {
				// If network server is a load balancer, getting server details fails due to missing permissions as load balancers underlying server instance is hidden from the user.
				serverRows = append(serverRows, output.TableRow{
					server.ServerUUID, server.ServerTitle, "", "",
				})
			} else {
				serverRows = append(serverRows, output.TableRow{
					fetched.UUID,
					fetched.Title,
					fetched.Hostname,
					fetched.State,
				})
			}
		}

		combined = append(combined, output.CombinedSection{
			Key:   "servers",
			Title: "Servers:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Header: "UUID", Key: "uuid", Colour: ui.DefaultUUUIDColours},
					{Header: "Title", Key: "title"},
					{Header: "Hostname", Key: "hostname", Format: format.PossiblyUnknownString},
					{Header: "State", Key: "state", Format: format.ServerState},
				},
				Rows: serverRows,
			},
		})
	}

	return output.MarshaledWithHumanOutput{
		Value:  network,
		Output: combined,
	}, nil
}

func formatShowDHCPDNS(val interface{}) (text.Colors, string, error) {
	addresses, ok := val.([]string)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse IP addresses from %T, expected []string", val)
	}

	return nil, ui.DefaultAddressColours.Sprint(strings.Join(addresses, " ")), nil
}
