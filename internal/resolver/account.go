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
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, account := range accounts {
			rv.AddMatch(account.Username, MatchTitle(arg, account.Username))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingAccount) PositionalArgumentHelp() string {
	return "<Username...>"
}
