package resolver

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

// CachingServerGroup implements resolver for servers, caching the results
type CachingServerGroup struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingServerGroup{}

// Get implements ResolutionProvider.Get
func (s CachingServerGroup) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	serverGroups, err := svc.GetServerGroups(ctx, &request.GetServerGroupsRequest{})
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, serverGroup := range serverGroups {
			if MatchArgWithWhitespace(arg, serverGroup.Title) || serverGroup.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = serverGroup.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingServerGroup) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
