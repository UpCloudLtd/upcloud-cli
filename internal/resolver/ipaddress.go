package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

// CachingIPAddress implements resolver for ip addresses that resolve with ptr records, caching the results
type CachingIPAddress struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingIPAddress{}

// Get implements ResolutionProvider.Get
func (s CachingIPAddress) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	ipaddresses, err := svc.GetIPAddresses(ctx)
	if err != nil {
		return nil, err
	}
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, ipAddress := range ipaddresses.IPAddresses {
			rv.AddMatch(ipAddress.Address, MatchTitle(arg, ipAddress.PTRRecord, ipAddress.Address))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingIPAddress) PositionalArgumentHelp() string {
	return "<Address/PTRRecord...>"
}
