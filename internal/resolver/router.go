package resolver

import (
	"context"
	"errors"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CachingRouter implements resolver for routers by uuid or name, caching the results
type CachingRouter struct {
	cached []upcloud.Router
}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = &CachingRouter{}

// Get implements ResolutionProvider.Get
func (s *CachingRouter) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	routers, err := svc.GetRouters(ctx)
	if err != nil {
		return nil, err
	}
	s.cached = routers.Routers
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, router := range s.cached {
			rv.AddMatch(router.UUID, MatchTitle(arg, router.Name))
			rv.AddMatch(router.UUID, MatchUUID(arg, router.UUID))
		}
		return rv
	}, nil
}

// GetCached is a helper method for commands to use when they need to get an item from the cached results
func (s *CachingRouter) GetCached(uuid string) (upcloud.Router, error) {
	if s.cached == nil {
		return upcloud.Router{}, errors.New("caching network does not have a cache initialized")
	}
	for _, router := range s.cached {
		if router.UUID == uuid {
			return router, nil
		}
	}
	return upcloud.Router{}, NotFoundError(uuid)
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *CachingRouter) PositionalArgumentHelp() string {
	return "<UUID/Name...>"
}
