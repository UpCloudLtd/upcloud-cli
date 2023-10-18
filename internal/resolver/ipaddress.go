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
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, ipAddress := range ipaddresses.IPAddresses {
			if ipAddress.PTRRecord == arg || ipAddress.Address == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = ipAddress.Address
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingIPAddress) PositionalArgumentHelp() string {
	return "<Address/PTRRecord...>"
}
