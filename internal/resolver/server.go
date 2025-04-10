package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CachingServer implements resolver for servers, caching the results
type CachingServer struct {
	Cache[upcloud.Server]
}

// make sure we implement the ResolutionProvider interfaces
var (
	_ ResolutionProvider                        = &CachingServer{}
	_ CachingResolutionProvider[upcloud.Server] = &CachingServer{}
)

// Get implements ResolutionProvider.Get
func (s *CachingServer) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	servers, err := svc.GetServers(ctx)
	if err != nil {
		return nil, err
	}

	for _, server := range servers.Servers {
		s.AddCached(server.UUID, server)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, server := range servers.Servers {
			rv.AddMatch(server.UUID, MatchTitle(arg, server.Title, server.Hostname))
			rv.AddMatch(server.UUID, MatchUUID(arg, server.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingServer) PositionalArgumentHelp() string {
	return "<UUID/Title/Hostname...>"
}
