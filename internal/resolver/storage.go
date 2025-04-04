package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CachingStorage implements resolver for storages, caching the results
type CachingStorage struct {
	Cache[upcloud.Storage]

	Access string
}

// make sure we implement the ResolutionProvider interfaces
var (
	_ ResolutionProvider                         = &CachingStorage{}
	_ CachingResolutionProvider[upcloud.Storage] = &CachingStorage{}
)

// Get implements ResolutionProvider.Get
func (s *CachingStorage) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	storages, err := svc.GetStorages(ctx, &request.GetStoragesRequest{Access: s.Access})
	if err != nil {
		return nil, err
	}

	for _, storage := range storages.Storages {
		s.AddCached(storage.UUID, storage)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, storage := range storages.Storages {
			rv.AddMatch(storage.UUID, MatchTitle(arg, storage.Title))
			rv.AddMatch(storage.UUID, MatchUUID(arg, storage.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *CachingStorage) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
