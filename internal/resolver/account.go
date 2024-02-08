package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
)

// CachingAccount implements resolver for servers, caching the results
type CachingAccount struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingAccount{}

// Get implements ResolutionProvider.Get
func (s CachingAccount) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	accounts, err := svc.GetAccountList(ctx)
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, account := range accounts {
			if MatchArgWithWhitespace(arg, account.Username) {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = account.Username
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingAccount) PositionalArgumentHelp() string {
	return "<Username...>"
}
