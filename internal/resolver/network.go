package resolver

import (
	"errors"
	internal "github.com/UpCloudLtd/cli/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
)

// CachingNetwork implements resolver for networks, caching the results
type CachingNetwork struct {
	cached []upcloud.Network
}

// make sure we implement the ResolutionProvider interface
var _ ResolutionProvider = &CachingNetwork{}

// Get implements ResolutionProvider.Get
func (s *CachingNetwork) Get(svc internal.AllServices) (Resolver, error) {
	networks, err := svc.GetNetworks()
	if err != nil {
		return nil, err
	}
	s.cached = networks.Networks
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, network := range s.cached {
			if network.Name == arg || network.UUID == arg {
				if rv != "" {
					return "", AmbiguousResolutionError(arg)
				}
				rv = network.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", NotFoundError(arg)
	}, nil
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
