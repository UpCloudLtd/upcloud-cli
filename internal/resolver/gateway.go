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
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, gtw := range gateways {
			if MatchArgWithWhitespace(arg, gtw.Name) || gtw.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = gtw.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingGateway) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
