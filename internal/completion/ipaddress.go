package completion

import (
	internal "github.com/UpCloudLtd/cli/internal/service"
	"github.com/spf13/cobra"
)

// IPAddress implements argument completion for ip addresses, by ptr record or the adddress itself
type IPAddress struct {
}

// Generate implements completion.Provider
func (s IPAddress) Generate(services internal.AllServices) (Completer, error) {
	ipAddresses, err := services.GetIPAddresses()
	if err != nil {
		return func(toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveDefault
		}, nil
	}
	return func(toComplete string) ([]string, cobra.ShellCompDirective) {
		var vals []string
		for _, v := range ipAddresses.IPAddresses {
			vals = append(vals, v.PTRRecord, v.Address)
		}
		return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
	}, nil
}
