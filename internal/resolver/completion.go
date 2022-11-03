package resolver

import (
	internal "github.com/UpCloudLtd/upcloud-cli/v2/internal/service"
)

// CompletionResolver implements resolver for servers, caching the results
type CompletionResolver struct{}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = CompletionResolver{}

// Get is just a passthrough to respect the lib
func (s CompletionResolver) Get(_ internal.AllServices) (Resolver, error) {
	return func(arg string) (uuid string, err error) {
		return arg, nil
	}, nil
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s CompletionResolver) PositionalArgumentHelp() string {
	return "<bash>"
}
