package completion

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/spf13/cobra"
)

// IPAddress implements argument completion for ip addresses, by ptr record or the adddress itself
type IPAddress struct {
}

// make sure IPAddress implements the interface
var _ Provider = IPAddress{}

// CompleteArgument implements completion.Provider
func (s IPAddress) CompleteArgument(svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	fmt.Println(svc.GetIPAddresses())
	ipAddresses, err := svc.GetIPAddresses()
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range ipAddresses.IPAddresses {
		vals = append(vals, v.PTRRecord, v.Address)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
