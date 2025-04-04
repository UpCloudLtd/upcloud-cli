package resolver

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

// CachingServerGroup implements resolver for servers, caching the results
type CachingServerGroup struct {
	Cache[upcloud.ServerGroup]
}

// make sure we implement the ResolutionProvider interface
var (
	_ ResolutionProvider                             = &CachingServerGroup{}
	_ CachingResolutionProvider[upcloud.ServerGroup] = &CachingServerGroup{}
)

// Get implements ResolutionProvider.Get
func (s *CachingServerGroup) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	serverGroups, err := svc.GetServerGroups(ctx, &request.GetServerGroupsRequest{})
	if err != nil {
		return nil, err
	}

	for _, serverGroup := range serverGroups {
		s.AddCached(serverGroup.UUID, serverGroup)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, serverGroup := range serverGroups {
			rv.AddMatch(serverGroup.UUID, MatchTitle(arg, serverGroup.Title))
			rv.AddMatch(serverGroup.UUID, MatchUUID(arg, serverGroup.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingServerGroup) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
