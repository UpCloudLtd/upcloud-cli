package namedargs

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
)

// ResolveNetwork resolves network UUID from values provided to named args (e.g., --network net-name)
func ResolveNetwork(exec commands.Executor, arg string) (string, error) {
	net, err := Resolve(&resolver.CachingNetwork{}, exec, arg)
	if err != nil {
		err = fmt.Errorf("could not resolve network: %w", err)
	}

	return net, err
}
