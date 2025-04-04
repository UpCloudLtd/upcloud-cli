package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CachingDatabase implements resolver for servers, caching the results
type CachingDatabase struct {
	Cache[upcloud.ManagedDatabase]
}

// make sure we implement the ResolutionProvider interface
var (
	_ ResolutionProvider                                 = &CachingDatabase{}
	_ CachingResolutionProvider[upcloud.ManagedDatabase] = &CachingDatabase{}
)

// Get implements ResolutionProvider.Get
func (s *CachingDatabase) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	databases, err := svc.GetManagedDatabases(ctx, &request.GetManagedDatabasesRequest{})
	if err != nil {
		return nil, err
	}
	for _, db := range databases {
		s.AddCached(db.UUID, db)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, db := range databases {
			rv.AddMatch(db.UUID, MatchTitle(arg, db.Title))
			rv.AddMatch(db.UUID, MatchUUID(arg, db.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingDatabase) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
