package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

// CachingServer implements resolver for servers, caching the results
type CachingServer struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingServer{}

// Get implements ResolutionProvider.Get
func (s CachingServer) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	servers, err := svc.GetServers(ctx)
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, server := range servers.Servers {
			if MatchArgWithWhitespace(arg, server.Title) || server.Hostname == arg || server.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = server.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingServer) PositionalArgumentHelp() string {
	return "<UUID/Title/Hostname...>"
}
