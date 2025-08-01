package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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

// KubernetesVersion implements argument completion for Kubernetes versions by version id.
type KubernetesVersion struct{}

// make sure Kubernetes implements the interface
var _ Provider = KubernetesVersion{}

// CompleteArgument implements completion.Provider
func (s KubernetesVersion) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	versions, err := svc.GetKubernetesVersions(ctx, &request.GetKubernetesVersionsRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, version := range versions {
		vals = append(vals, version.Id)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

// KubernetesPlan implements argument completion for Kubernetes plans.
type KubernetesPlan struct{}

// make sure KubernetesPlan implements the interface.
var _ Provider = KubernetesPlan{}

// CompleteArgument implements completion.Provider.
func (s KubernetesPlan) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	plans, err := svc.GetKubernetesPlans(ctx, &request.GetKubernetesPlansRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, plan := range plans {
		if !plan.Deprecated {
			vals = append(vals, plan.Name)
		}
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
