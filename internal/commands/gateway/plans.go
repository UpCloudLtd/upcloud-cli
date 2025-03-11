package gateway

import (
	"sort"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

// PlansCommand creates the "gateway plans" command
func PlansCommand() commands.Command {
	return &plansCommand{
		BaseCommand: commands.New("plans", "List gateway plans", "upctl gateway plans"),
	}
}

type plansCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *plansCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	plans, err := exec.All().GetGatewayPlans(exec.Context())
	if err != nil {
		return nil, err
	}

	sort.Slice(plans, func(i, j int) bool {
		if plans[i].ServerNumber != plans[j].ServerNumber {
			return plans[i].ServerNumber < plans[j].ServerNumber
		}

		if plans[i].PerGatewayBandwidthMbps != plans[j].PerGatewayBandwidthMbps {
			return plans[i].PerGatewayBandwidthMbps < plans[j].PerGatewayBandwidthMbps
		}

		if plans[i].PerGatewayMaxConnections != plans[j].PerGatewayMaxConnections {
			return plans[i].PerGatewayMaxConnections < plans[j].PerGatewayMaxConnections
		}

		return plans[i].VPNTunnelAmount < plans[j].VPNTunnelAmount
	})

	rows := []output.TableRow{}
	for _, p := range plans {
		rows = append(rows, output.TableRow{
			p.Name,
			p.ServerNumber,
			p.SupportedFeatures,
			p.PerGatewayBandwidthMbps,
			p.PerGatewayMaxConnections,
			p.VPNTunnelAmount,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "server_number", Header: "Server count"},
				{Key: "supported_features", Header: "Features", Format: format.StringSliceAnd},
				{Key: "per_gateway_bandwidth_mbps", Header: "Max bandwidth (Mbps)"},
				{Key: "per_gateway_max_connections", Header: "Max connections"},
				{Key: "vpn_tunnel_amount", Header: "VPN tunnel amount"},
			},
			Rows: rows,
		},
	}, nil
}
