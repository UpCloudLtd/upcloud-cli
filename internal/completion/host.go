package completion

import (
	"context"
	"strconv"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// HostID implements argument completion for host IDs.
type HostID struct{}

// make sure Token implements the interface
var _ Provider = HostID{}

// CompleteArgument implements completion.Provider
func (s HostID) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	hosts, err := svc.GetHosts(ctx)
	if err != nil {
		return None(toComplete)
	}
	vals := make([]string, 0, len(hosts.Hosts))
	for _, h := range hosts.Hosts {
		vals = append(vals, strconv.Itoa(h.ID))
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
