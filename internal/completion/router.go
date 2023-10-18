package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// Router implements argument completion for routers, by name or uuid.
type Router struct{}

// make sure Router implements the interface
var _ Provider = Router{}

// CompleteArgument implements completion.Provider
func (s Router) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	routers, err := svc.GetRouters(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range routers.Routers {
		vals = append(vals, v.UUID, v.Name)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
