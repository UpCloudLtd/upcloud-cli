package resolver

import "github.com/UpCloudLtd/cli/internal/service"

// Resolver represents the most basic argument resolver, a function that accepts and argument and returns an uuid (or error)
type Resolver func(arg string) (uuid string, err error)

// ResolutionProvider is an interface for commands that provide resolution, either custom or the built-in ones
type ResolutionProvider interface {
	Get(svc service.AllServices) (Resolver, error)
	PositionalArgumentHelp() string
}
