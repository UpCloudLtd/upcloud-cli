package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// ServerGroup implements argument completion for server groups, by uuid or title.
type ServerGroup struct{}

// make sure ServerGroup implements the interface
var _ Provider = ServerGroup{}

// CompleteArgument implements completion.Provider
func (s ServerGroup) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	vals, err := serverGroupCompletions(ctx, svc, true)
	if err != nil {
		return None(toComplete)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

// ServerGroupUUID implements argument completion for server group UUIDs.
type ServerGroupUUID struct{}

// make sure ServerGroupUUID implements the interface.
var _ Provider = ServerGroupUUID{}

// CompleteArgument implements completion.Provider
func (s ServerGroupUUID) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	vals, err := serverGroupCompletions(ctx, svc, false)
	if err != nil {
		return None(toComplete)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

func serverGroupCompletions(ctx context.Context, svc service.AllServices, withTitles bool) ([]string, error) {
	serverGroups, err := svc.GetServerGroups(ctx, &request.GetServerGroupsRequest{})
	if err != nil {
		return nil, err
	}
	var vals []string
	for _, v := range serverGroups {
		vals = append(vals, v.UUID)
		if withTitles {
			vals = append(vals, v.Title)
		}
	}
	return vals, nil
}
