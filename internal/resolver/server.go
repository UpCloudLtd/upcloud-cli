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
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, server := range servers.Servers {
			rv.AddMatch(server.UUID, MatchArgWithWhitespace(arg, server.Title))
			rv.AddMatch(server.UUID, MatchArgWithWhitespace(arg, server.Hostname))
			rv.AddMatch(server.UUID, MatchUUID(arg, server.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingServer) PositionalArgumentHelp() string {
	return "<UUID/Title/Hostname...>"
}
