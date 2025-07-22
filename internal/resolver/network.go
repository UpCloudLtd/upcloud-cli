package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CachingNetwork implements resolver for networks, caching the results
type CachingNetwork struct {
	Cache[upcloud.Network]
}

var (
	_ ResolutionProvider                         = &CachingNetwork{}
	_ CachingResolutionProvider[upcloud.Network] = &CachingNetwork{}
)

// Get implements ResolutionProvider.Get
func (s *CachingNetwork) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	networks, err := svc.GetNetworks(ctx)
	if err != nil {
		return nil, err
	}

	for _, network := range networks.Networks {
		s.AddCached(network.UUID, network)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, network := range networks.Networks {
			rv.AddMatch(network.UUID, MatchTitle(arg, network.Name))
			rv.AddMatch(network.UUID, MatchUUID(arg, network.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *CachingNetwork) PositionalArgumentHelp() string {
	return helpUUIDName
}
