package resolver

import (
	"context"
	"errors"

	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CachingNetwork implements resolver for networks, caching the results
type CachingNetwork struct {
	cached []upcloud.Network
}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = &CachingNetwork{}

func networkMatcher(cached []upcloud.Network) func(arg string) Resolved {
	return func(arg string) Resolved {
		rv := Resolved{Arg: arg}
		for _, network := range cached {
			rv.AddMatch(network.UUID, MatchTitle(arg, network.Name))
			rv.AddMatch(network.UUID, MatchUUID(arg, network.UUID))
		}
		return rv
	}
}

// Get implements ResolutionProvider.Get
func (s *CachingNetwork) Get(ctx context.Context, svc internal.AllServices) (Resolver, error) {
	networks, err := svc.GetNetworks(ctx)
	if err != nil {
		return nil, err
	}
	s.cached = networks.Networks
	return networkMatcher(s.cached), nil
}

// GetCached is a helper method for commands to use when they need to get an item from the cached results
func (s *CachingNetwork) GetCached(uuid string) (upcloud.Network, error) {
	if s.cached == nil {
		return upcloud.Network{}, errors.New("caching network does not have a cache initialized")
	}
	for _, network := range s.cached {
		if network.UUID == uuid {
			return network, nil
		}
	}
	return upcloud.Network{}, NotFoundError(uuid)
}

// Resolve is a helper method for commands to resolve networks inside Execute(), outside arguments
func (s *CachingNetwork) Resolve(arg string) (resolved string, err error) {
	if s.cached == nil {
		return "", errors.New("caching network does not have a cache initialized")
	}

	r := networkMatcher(s.cached)(arg)
	return r.GetOnly()
}

// PositionalArgumentHelp implements resolver.ResolutionProvider
func (s *CachingNetwork) PositionalArgumentHelp() string {
	return "<UUID/Name...>"
}
