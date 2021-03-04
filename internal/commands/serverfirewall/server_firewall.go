package serverfirewall

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/cobra"
)

const (
	PositionalArgHelp = "/server/<UUID/firewall_rule>"
)

// CachedServers stores a cached list of servers fetched from the service
// TODO: remove the cross-command dependencies
var CachedServers []upcloud.Server

func BaseServerFirewallCommand() commands.Command {
	return &serverFirewallCommand{commands.New("firewall", "Manage server firewall rules. Enabling or disabling the firewall is done in server modify.")}
}

type serverFirewallCommand struct {
	*commands.BaseCommand
}

func matchServers(servers []upcloud.Server, searchVal string) []*upcloud.Server {
	var r []*upcloud.Server
	for _, server := range servers {
		server := server
		if server.Title == searchVal || server.Hostname == searchVal || server.UUID == searchVal {
			r = append(r, &server)
		}
	}
	return r
}

// SearchSingleServer returns exactly single server with uuid, name or title matching term.
// The service will be hit once to get the list of all servers and results will be cached afterwards.
func SearchSingleServer(term string, service service.Server) (*upcloud.Server, error) {
	servers, err := searchServer(&CachedServers, service, term, true)
	if err != nil {
		return nil, err
	}
	return servers[0], nil
}

func searchServer(serversPtr *[]upcloud.Server, service service.Server, uuidOrHostnameOrTitle string, unique bool) ([]*upcloud.Server, error) {
	if serversPtr == nil || service == nil {
		return nil, fmt.Errorf("no servers or service passed")
	}
	servers := *serversPtr
	if len(CachedServers) == 0 {
		res, err := service.GetServers()
		if err != nil {
			return nil, err
		}
		servers = res.Servers
		*serversPtr = servers
	}
	matched := matchServers(servers, uuidOrHostnameOrTitle)
	if len(matched) == 0 {
		return nil, fmt.Errorf("no server with uuid, name or title %q was found", uuidOrHostnameOrTitle)
	}
	if len(matched) > 1 && unique {
		return nil, fmt.Errorf("multiple servers matched to query %q, use UUID to specify", uuidOrHostnameOrTitle)
	}
	return matched, nil
}

func searchAllServers(terms []string, service service.Server, unique bool) ([]string, error) {
	return commands.SearchResources(
		terms,
		func(term string) (interface{}, error) {
			return searchServer(&CachedServers, service, term, unique)
		},
		func(in interface{}) string { return in.(*upcloud.Server).UUID })
}

// GetServerArgumentCompletionFunction returns a bash completion function for servers
func GetServerArgumentCompletionFunction(s service.Server) func(toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(toComplete string) ([]string, cobra.ShellCompDirective) {
		servers, err := s.GetServers()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range servers.Servers {
			vals = append(vals, v.UUID, v.Hostname, v.Title)
		}
		return commands.MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
	}
}
