package resolver

import (
	"errors"
	internal "github.com/UpCloudLtd/cli/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
)

// CachingRouter implements resolver for routers by uuid or name, caching the results
type CachingRouter struct {
	cached []upcloud.Router
}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = &CachingRouter{}

// Get implements ResolutionProvider.Get
func (s *CachingRouter) Get(svc internal.AllServices) (Resolver, error) {
	routers, err := svc.GetRouters()
	if err != nil {
		return nil, err
	}
	s.cached = routers.Routers
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, router := range s.cached {
			if router.Name == arg || router.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = router.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
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
