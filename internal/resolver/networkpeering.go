package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CachingNetworkPeering resolver for network peerings, caching the results
type CachingNetworkPeering struct {
	Cache[upcloud.NetworkPeering]
}

var (
	_ ResolutionProvider                                = &CachingNetworkPeering{}
	_ CachingResolutionProvider[upcloud.NetworkPeering] = &CachingNetworkPeering{}
)

// Get implements ResolutionProvider.Get
func (s *CachingNetworkPeering) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	peerings, err := svc.GetNetworkPeerings(ctx)
	if err != nil {
		return nil, err
	}

	for _, peering := range peerings {
		s.AddCached(peering.UUID, peering)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, peering := range peerings {
			s.AddCached(peering.UUID, peering)

			rv.AddMatch(peering.UUID, MatchTitle(arg, peering.Name))
			rv.AddMatch(peering.UUID, MatchUUID(arg, peering.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingNetworkPeering) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
