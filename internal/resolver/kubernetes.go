package resolver

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// CachingKubernetes implements resolver for Kubernetes clusters, caching the results
type CachingKubernetes struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingKubernetes{}

// Get implements ResolutionProvider.Get
func (s CachingKubernetes) Get(ctx context.Context, svc service.AllServices) (Resolver, error) {
	clusters, err := svc.GetKubernetesClusters(ctx, &request.GetKubernetesClustersRequest{})
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, cluster := range clusters {
			if cluster.Name == arg || cluster.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = cluster.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingKubernetes) PositionalArgumentHelp() string {
	return "<UUID/Name...>" //nolint:goconst
}
