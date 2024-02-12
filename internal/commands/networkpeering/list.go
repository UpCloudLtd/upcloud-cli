package networkpeering

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
)

// ListCommand creates the "networkpeering list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List network peerings", "upctl networkpeering list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	peerings, err := svc.GetNetworkPeerings(exec.Context())
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, peering := range peerings {
		peerNetwork := ""
		if len(peering.PeerNetwork.IPNetworks) > 0 {
			peerNetwork = peering.PeerNetwork.IPNetworks[0].Address
		}

		rows = append(rows, output.TableRow{
			peering.UUID,
			peering.Name,
			peering.Network.IPNetworks[0].Address,
			peerNetwork,
			peering.State,
		})
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: peerings,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "network", Header: "Network", Colour: ui.DefaultAddressColours},
				{Key: "peer_network", Header: "Peer Network", Colour: ui.DefaultAddressColours},
				{Key: "status", Header: "Status"},
			},
			Rows: rows,
		},
	}, nil
}
