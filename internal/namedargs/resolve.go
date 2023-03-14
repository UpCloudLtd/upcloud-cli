package namedargs

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
)

// Resolve initializes given resolution provider and uses it to resolve given argument
func Resolve(provider resolver.ResolutionProvider, exec commands.Executor, arg string) (string, error) {
	resolver, err := provider.Get(exec.Context(), exec.All())
	if err != nil {
		return "", fmt.Errorf("could not initialize resolver: %w", err)
	}

	return resolver(arg)
}

// ResolveNetwork resolves network UUID from values provided to named args (e.g., --network net-name)
func ResolveNetwork(exec commands.Executor, arg string) (string, error) {
	net, err := Resolve(&resolver.CachingNetwork{}, exec, arg)
	if err != nil {
		err = fmt.Errorf("could not resolve network: %w", err)
	}
	return net, err
}
