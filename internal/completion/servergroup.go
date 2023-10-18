package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// ServerGroup implements argument completion for server groups, by uuid or title.
type ServerGroup struct{}

// make sure ServerGroup implements the interface
var _ Provider = ServerGroup{}

// CompleteArgument implements completion.Provider
func (s ServerGroup) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	serverGroups, err := svc.GetServerGroups(ctx, &request.GetServerGroupsRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range serverGroups {
		vals = append(vals, v.UUID, v.Title)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
