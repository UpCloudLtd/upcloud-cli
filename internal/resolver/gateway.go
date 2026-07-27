package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CachingGatewayimplements resolver for gateways, caching the results
type CachingGateway struct {
	Cache[upcloud.Gateway]
}

// make sure we implement the ResolutionProvider interfaces
var (
	_ ResolutionProvider                         = &CachingGateway{}
	_ CachingResolutionProvider[upcloud.Gateway] = &CachingGateway{}
)

// Get implements ResolutionProvider.Get
func (s *CachingGateway) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	gateways, err := svc.GetGateways(ctx)
	if err != nil {
		return nil, err
	}
	for _, gw := range gateways {
		s.AddCached(gw.UUID, gw)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, gtw := range gateways {
			rv.AddMatch(gtw.UUID, MatchTitle(arg, gtw.Name))
			rv.AddMatch(gtw.UUID, MatchUUID(arg, gtw.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingGateway) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
