package resolver

import (
	"context"
	"errors"

	internal "github.com/UpCloudLtd/upcloud-cli/v2/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// CachingStorage implements resolver for storages, caching the results
type CachingStorage struct {
	cachedStorages *upcloud.Storages
}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = &CachingStorage{}

func storageMatcher(cached []upcloud.Storage) func(arg string) (uuid string, err error) {
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, storage := range cached {
			if MatchArgWithWhitespace(arg, storage.Title) || storage.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = storage.UUID
			}
		}

		if rv != "" {
			return rv, nil
		}

		return "", NotFoundError(arg)
	}
}

// Get implements ResolutionProvider.Get
func (s *CachingStorage) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	var err error
	s.cachedStorages, err = svc.GetStorages(ctx, &request.GetStoragesRequest{})
	if err != nil {
		return nil, err
	}
	return storageMatcher(s.cachedStorages.Storages), nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *CachingStorage) PositionalArgumentHelp() string {
	return "<UUID/Title...>"
}

// Resolve is a helper method for commands to resolve networks inside Execute(), outside arguments
func (s *CachingStorage) Resolve(arg string) (resolved string, err error) {
	if s.cachedStorages == nil {
		return "", errors.New("caching storage does not have a cache initialized")
	}

	return storageMatcher(s.cachedStorages.Storages)(arg)
}

// GetCached is a helper method for commands to use when they need to get an item from the cached results
func (s *CachingStorage) GetCached(uuid string) (upcloud.Storage, error) {
	if s.cachedStorages == nil {
		return upcloud.Storage{}, errors.New("caching storages does not have a cache initialized")
	}
	for _, storage := range s.cachedStorages.Storages {
		if storage.UUID == uuid {
			return storage, nil
		}
	}
	return upcloud.Storage{}, NotFoundError(uuid)
}
