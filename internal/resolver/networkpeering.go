package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

// CachingNetworkPeering resolver for network peerings, caching the results
type CachingNetworkPeering struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingNetworkPeering{}

// Get implements ResolutionProvider.Get
func (s CachingNetworkPeering) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	gateways, err := svc.GetNetworkPeerings(ctx)
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, peering := range gateways {
			if MatchArgWithWhitespace(arg, peering.Name) || peering.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = peering.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingNetworkPeering) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
