package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CachingObjectStorage implements resolver for ObjectStorages, caching the results
type CachingObjectStorage struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingObjectStorage{}

// Get implements ResolutionProvider.Get
func (s CachingObjectStorage) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	objectstorages, err := svc.GetManagedObjectStorages(ctx, &request.GetManagedObjectStoragesRequest{})
	if err != nil {
		return nil, err
	}
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, objsto := range objectstorages {
			rv.AddMatch(objsto.UUID, MatchTitle(arg, objsto.Name))
			rv.AddMatch(objsto.UUID, MatchUUID(arg, objsto.UUID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingObjectStorage) PositionalArgumentHelp() string {
	return helpUUIDTitle
}
