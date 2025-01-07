package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

// CachingGatewayimplements resolver for gateways, caching the results
type CachingGateway struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingGateway{}

// Get implements ResolutionProvider.Get
func (s CachingGateway) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	gateways, err := svc.GetGateways(ctx)
	if err != nil {
		return nil, err
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
