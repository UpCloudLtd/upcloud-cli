package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CachingLoadBalancer implements resolver for servers, caching the results
type CachingLoadBalancer struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingLoadBalancer{}

// Get implements ResolutionProvider.Get
func (s CachingLoadBalancer) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	loadbalancers, err := svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{})
	if err != nil {
		return nil, err
	}
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, lb := range loadbalancers {
			rv.AddMatch(lb.UUID, MatchTitle(arg, lb.Name))
			rv.AddMatch(lb.UUID, MatchUUID(arg, lb.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingLoadBalancer) PositionalArgumentHelp() string {
	return "<UUID/Name...>"
}
