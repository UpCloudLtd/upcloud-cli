package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/service"
	"github.com/spf13/cobra"
)

// Server implements argument completion for servers, by uuid, name or hostname.
type Server struct{}

// make sure Server implements the interface
var _ Provider = Server{}

// CompleteArgument implements completion.Provider
func (s Server) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	servers, err := svc.GetServers(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range servers.Servers {
		vals = append(vals, v.UUID, v.Hostname, v.Title)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
