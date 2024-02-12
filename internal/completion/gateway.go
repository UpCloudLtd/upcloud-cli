package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// Gateway implements argument completion for gateways, by uuid or name.
type Gateway struct{}

// make sure Gateway implements the interface
var _ Provider = Gateway{}

// CompleteArgument implements completion.Provider
func (s Gateway) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	gateways, err := svc.GetGateways(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, gtw := range gateways {
		vals = append(vals, gtw.UUID, gtw.Name)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
