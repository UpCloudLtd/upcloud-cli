package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// NetworkPeering implements argument completion for network peerings, by uuid or name.
type NetworkPeering struct{}

// make sure NetworkPeering implements the interface
var _ Provider = NetworkPeering{}

// CompleteArgument implements completion.Provider
func (s NetworkPeering) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	peerings, err := svc.GetNetworkPeerings(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, peering := range peerings {
		vals = append(vals, peering.UUID, peering.Name)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
