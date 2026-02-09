package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CachingFileStorage implements resolver for FileStorages, caching the results
type CachingFileStorage struct {
	Cache[upcloud.FileStorage]
}

// make sure we implement the ResolutionProvider interfaces
var (
	_ ResolutionProvider                             = &CachingFileStorage{}
	_ CachingResolutionProvider[upcloud.FileStorage] = &CachingFileStorage{}
)

// Get implements ResolutionProvider.Get
func (s *CachingFileStorage) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	filestorages, err := svc.GetFileStorages(ctx, &request.GetFileStoragesRequest{})
	if err != nil {
		return nil, err
	}
	for _, filesto := range filestorages {
		s.AddCached(filesto.UUID, filesto)
	}

	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, filesto := range filestorages {
			rv.AddMatch(filesto.UUID, MatchTitle(arg, filesto.Name))
			rv.AddMatch(filesto.UUID, MatchUUID(arg, filesto.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingFileStorage) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
