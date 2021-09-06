package completion

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/service"
	"github.com/spf13/cobra"
)

// Server implements argument completion for routers, by uuid, name or hostname.
type Server struct {
}

// make sure Server implements the interface.
var _ Provider = Server{}

// CompleteArgument implements completion.Provider.
func (s Server) CompleteArgument(svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	servers, err := svc.GetServers()
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range servers.Servers {
		vals = append(vals, v.UUID, v.Hostname, v.Title)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
