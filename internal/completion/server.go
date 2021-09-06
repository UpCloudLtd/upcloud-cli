package completion

import (
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/upcloud-cli/internal/service"
)

// Server implements argument completion for routers, by uuid, name or hostname.
type Server struct{}

// make sure Server implements the interface.
var _ Provider = Server{}

// CompleteArgument implements completion.Provider.
func (s Server) CompleteArgument(svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	servers, err := svc.GetServers()
	if err != nil {
		return None(toComplete)
	}
	vals := make([]string, 0, len(servers.Servers))
	for _, v := range servers.Servers {
		vals = append(vals, v.UUID, v.Hostname, v.Title)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
