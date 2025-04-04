package resolver

import (
	"context"
	"errors"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CachingStorage implements resolver for storages, caching the results
type CachingStorage struct {
	Cache[upcloud.Storage]
}

// make sure we implement the ResolutionProvider interfaces
var (
	_ ResolutionProvider                         = &CachingStorage{}
	_ CachingResolutionProvider[upcloud.Storage] = &CachingStorage{}
)

func (s *CachingStorage) matchCached() func(arg string) Resolved {
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, storage := range s.cache {
			rv.AddMatch(storage.UUID, MatchTitle(arg, storage.Title))
			rv.AddMatch(storage.UUID, MatchUUID(arg, storage.UUID))
		}
		return rv
	}
}

// Get implements ResolutionProvider.Get
func (s *CachingStorage) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	storages, err := svc.GetStorages(ctx, &request.GetStoragesRequest{})
	if err != nil {
		return nil, err
	}

	for _, storage := range storages.Storages {
		s.AddCached(storage.UUID, storage)
	}

	return s.matchCached(), nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *CachingStorage) PositionalArgumentHelp() string {
	return helpUUIDTitle
}

// Resolve is a helper method for commands to resolve networks inside Execute(), outside arguments
func (s *CachingStorage) Resolve(arg string) (resolved string, err error) {
	if s.cache == nil {
		return "", errors.New("caching storage does not have a cache initialized")
	}

	r := s.matchCached()(arg)
	return r.GetOnly()
}
