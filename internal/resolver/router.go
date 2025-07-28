package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CachingRouter implements resolver for routers by uuid or name, caching the results
type CachingRouter struct {
	Cache[upcloud.Router]

	Type string
}

var (
	_ ResolutionProvider                        = &CachingRouter{}
	_ CachingResolutionProvider[upcloud.Router] = &CachingRouter{}
)

// Get implements ResolutionProvider.Get
func (s *CachingRouter) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	routers, err := svc.GetRouters(ctx)
	if err != nil {
		return nil, err
	}

	for _, router := range routers.Routers {
		if s.Type != "" && router.Type != s.Type {
			continue
		}
		s.AddCached(router.UUID, router)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, router := range s.cache {
			rv.AddMatch(router.UUID, MatchTitle(arg, router.Name))
			rv.AddMatch(router.UUID, MatchUUID(arg, router.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *CachingRouter) PositionalArgumentHelp() string {
	return helpUUIDName
}
