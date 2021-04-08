package completion

import (
	internal "github.com/UpCloudLtd/cli/internal/service"
	"github.com/spf13/cobra"
)

// Network implements argument completion for networks, by name or uuid.
type Network struct {
}

// Generate implements completion.Provider
func (n Network) Generate(services internal.AllServices) (Completer, error) {
	networks, err := services.GetNetworks()
	if err != nil {
		return func(toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveDefault
		}, nil
	}
	return func(toComplete string) ([]string, cobra.ShellCompDirective) {
		var vals []string
		for _, v := range networks.Networks {
			vals = append(vals, v.UUID, v.Name)
		}
		return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
	}, nil
}
