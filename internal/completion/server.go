package completion

import (
	"context"
	"slices"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// Server implements argument completion for servers, by uuid, name or hostname.
type Server struct{}

// make sure Server implements the interface
var _ Provider = Server{}

// CompleteArgument implements completion.Provider
func (s Server) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeServers(ctx, svc, toComplete)
}

// StartedServer implements argument completion for started servers, by uuid, name or hostname.
type StartedServer struct{}

// make sure StartedServer implements the interface
var _ Provider = StartedServer{}

// CompleteArgument implements completion.Provider
func (s StartedServer) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeServers(ctx, svc, toComplete, "started")
}

// Stopped implements argument completion for stopped servers, by uuid, name or hostname.
type StoppedServer struct{}

// make sure StoppedServer implements the interface
var _ Provider = StoppedServer{}

// CompleteArgument implements completion.Provider
func (s StoppedServer) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeServers(ctx, svc, toComplete, "stopped")
}

func completeServers(ctx context.Context, svc service.AllServices, toComplete string, states ...string) ([]string, cobra.ShellCompDirective) {
	servers, err := svc.GetServers(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range servers.Servers {
		if len(states) == 0 || slices.Contains(states, v.State) {
			vals = append(vals, v.UUID, v.Hostname, v.Title)
		}
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
