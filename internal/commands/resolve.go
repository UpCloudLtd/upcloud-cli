package commands

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
)

func Resolve(provider resolver.ResolutionProvider, exec Executor, arg string) (string, error) {
	resolver, err := provider.Get(exec.Context(), exec.All())
	if err != nil {
		return "", fmt.Errorf("could not initialize resolver: %w", err)
	}

	return resolver(arg)
}

func ResolveNetwork(exec Executor, arg string) (string, error) {
	net, err := Resolve(&resolver.CachingNetwork{}, exec, arg)
	if err != nil {
		err = fmt.Errorf("could not resolve network: %w", err)
	}
	return net, err
}
