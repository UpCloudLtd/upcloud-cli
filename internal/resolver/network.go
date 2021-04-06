package resolver

import (
	"fmt"
	internal "github.com/UpCloudLtd/cli/internal/service"
)

// CachingNetwork implements resolver for networks, caching the results
type CachingNetwork struct{}

var _ ResolutionProvider = CachingNetwork{}

// Get implements ResolutionProvider.Get
func (s CachingNetwork) Get(svc internal.AllServices) (Resolver, error) {
	networks, err := svc.GetNetworks()
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, network := range networks.Networks {
			if network.Name == arg || network.UUID == arg {
				if rv != "" {
					return "", fmt.Errorf("'%v' is ambiguous, found multiple networks matching", arg)
				}
				rv = network.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", fmt.Errorf("no network found matching '%v'", arg)
	}, nil
}
