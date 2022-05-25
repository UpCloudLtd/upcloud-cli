package resolver

import (
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// CachingDatabase implements resolver for servers, caching the results
type CachingDatabase struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingDatabase{}

// Get implements ResolutionProvider.Get
func (s CachingDatabase) Get(svc internal.AllServices) (Resolver, error) {
	databases, err := svc.GetManagedDatabases(&request.GetManagedDatabasesRequest{})
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, db := range databases {
			if db.Title == arg || db.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = db.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingDatabase) PositionalArgumentHelp() string {
	return "<UUID/Title...>"
}
