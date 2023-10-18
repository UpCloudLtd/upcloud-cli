package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/spf13/cobra"
)

// Network implements argument completion for networks, by name or uuid.
type Network struct{}

// make sure Network implements the interface
var _ Provider = Network{}

// CompleteArgument implements completion.Provider
func (s Network) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	networks, err := svc.GetNetworks(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range networks.Networks {
		// Only add completions for private networks
		if v.Type == "private" {
			vals = append(vals, v.UUID, v.Name)
		}
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
