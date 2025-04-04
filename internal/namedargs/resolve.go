package namedargs

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
)

// Resolve initializes given resolution provider and uses it to resolve given argument
func Resolve(provider resolver.ResolutionProvider, exec commands.Executor, arg string) (string, error) {
	resolver, err := provider.Get(exec.Context(), exec.All())
	if err != nil {
		return "", fmt.Errorf("could not initialize resolver: %w", err)
	}

	resolved := resolver(arg)
	return resolved.GetOnly()
}

// Resolve initializes given resolution provider, uses it to resolve given argument, and returns the cached resource.
func Get[T any](provider resolver.CachingResolutionProvider[T], exec commands.Executor, arg string) (T, error) {
	resolver, err := provider.Get(exec.Context(), exec.All())
	if err != nil {
		return *new(T), fmt.Errorf("could not initialize resolver: %w", err)
	}

	resolved := resolver(arg)
	uuid, err := resolved.GetOnly()
	if err != nil {
		return *new(T), err
	}
	return provider.GetCached(uuid)
}
