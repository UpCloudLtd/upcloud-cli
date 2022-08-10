package completion

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/spf13/cobra"
)

// Network implements argument completion for networks, by name or uuid.
type Network struct{}

// make sure Network implements the interface
var _ Provider = Network{}

// CompleteArgument implements completion.Provider
func (s Network) CompleteArgument(svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	networks, err := svc.GetNetworks()
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	// XXX: filter networks as it include all public/private prefixes
	for _, v := range networks.Networks {
		vals = append(vals, v.UUID, v.Name)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
