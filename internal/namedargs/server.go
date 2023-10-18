package namedargs

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
)

// ResolveServer resolves server UUID from values provided to named args (e.g., --server server-name)
func ResolveServer(exec commands.Executor, arg string) (string, error) {
	net, err := Resolve(&resolver.CachingServer{}, exec, arg)
	if err != nil {
		err = fmt.Errorf("could not resolve server: %w", err)
	}

	return net, err
}
