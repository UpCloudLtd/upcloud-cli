package server

import (
	"fmt"
	"time"

	// "github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

const (
	minStorageSize   = 10
	maxServerActions = 10
	// PositionalArgHelp is the helper string for server arguments
	// TODO: remove the cross-command dependencies
	PositionalArgHelp = "<UUID/Title/Hostname...>"

	//Server state related consts
	defaultStopType             = request.ServerStopTypeSoft
	defaultRestartTimeout       = time.Duration(120) * time.Second
	defaultRestartTimeoutAction = request.RestartTimeoutActionIgnore
)

// CachedServers stores a cached list of servers fetched from the service
// TODO: remove the cross-command dependencies
var CachedServers []upcloud.Server

// BaseServerCommand crestes the base "server" command
func BaseServerCommand() commands.Command {
	return &serverCommand{
		Command: commands.New("server", "Manage servers"),
	}
}

type serverCommand struct {
	commands.Command
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

// SearchAllServers searches for servers based on terms from service
// TODO: refactor?
func SearchAllServers(terms []string, service service.Server, unique bool) ([]string, error) {
	return commands.SearchResources(
		terms,
		func(term string) (interface{}, error) {
			return searchServer(&CachedServers, service, term, unique)
		},
		func(in interface{}) string { return in.(*upcloud.Server).UUID })
}

// Request represents a request for all servers
type Request struct {
	ExactlyOne   bool
	BuildRequest func(server string) interface{}
	Service      service.Server
	Handler      ui.Handler
}

// Send searches for all servers and calls Request.BuildRequest on them
func (s Request) Send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single server uuid is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one server uuid is required")
	}

	servers, err := SearchAllServers(args, s.Service, true)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, server := range servers {
		requests = append(requests, s.BuildRequest(server))
	}

	return s.Handler.Handle(requests)
}

// GetServerArgumentCompletionFunction returns a bash completion function for servers
func GetServerArgumentCompletionFunction(conf *config.Config) func(toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(toComplete string) ([]string, cobra.ShellCompDirective) {
		svc := conf.Service.Server()
		servers, err := svc.GetServers()
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

func serverStateWaiter(uuid, state, msg string, service service.Server, logline *ui.LogEntry) func() error {
	return func() error {
		for {
			time.Sleep(100 * time.Millisecond)
			details, err := service.GetServerDetails(&request.GetServerDetailsRequest{UUID: uuid})
			if err != nil {
				return err
			}
			if details.State == state {
				return nil
			}
			logline.SetMessage(fmt.Sprintf("%s: waiting to start (%v)", msg, details.State))
		}
	}
}
