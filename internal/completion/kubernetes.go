package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/cobra"
)

// Kubernetes implements argument completion for Kubernetes clusters, by uuid or name.
type Kubernetes struct{}

// make sure Kubernetes implements the interface
var _ Provider = Kubernetes{}

// CompleteArgument implements completion.Provider
func (s Kubernetes) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	clusters, err := svc.GetKubernetesClusters(ctx, &request.GetKubernetesClustersRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, cluster := range clusters {
		vals = append(vals, cluster.UUID, cluster.Name)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
