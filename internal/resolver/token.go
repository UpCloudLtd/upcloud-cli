package resolver

import (
	"context"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CachingToken implements resolver for tokens, caching the results.
type CachingToken struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CachingToken{}

// Get implements ResolutionProvider.Get
func (s CachingToken) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	tokens, err := svc.GetTokens(ctx, &request.GetTokensRequest{})
	if err != nil {
		return nil, err
	}
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, token := range *tokens {
			rv.AddMatch(token.ID, MatchTitle(arg, token.Name))
			rv.AddMatch(token.ID, MatchUUID(arg, token.ID))
		}
		return rv
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CachingToken) PositionalArgumentHelp() string {
	return "<ID...>"
}
