package resolver

import (
	internal "github.com/UpCloudLtd/cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

// CachingServer implements resolver for servers, caching the results
type CachingStorage struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingStorage{}

// Get implements ResolutionProvider.Get
func (s CachingStorage) Get(svc internal.AllServices) (Resolver, error) {
	storages, err := svc.GetStorages(&request.GetStoragesRequest{})
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, storage := range storages.Storages {
			if storage.Title == arg || storage.UUID == arg {
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
	}, nil
}
