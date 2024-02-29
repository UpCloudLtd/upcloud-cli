package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
)

// LoadBalancer implements argument completion for load balancers, by uuid or name.
type LoadBalancer struct{}

// make sure LoadBalancer implements the interface
var _ Provider = LoadBalancer{}

// CompleteArgument implements completion.Provider
func (s LoadBalancer) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	loadbalancers, err := svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, lb := range loadbalancers {
		vals = append(vals, lb.UUID, lb.Name)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
