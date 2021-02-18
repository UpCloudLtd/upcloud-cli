package server

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/spf13/cobra"
	"time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

const (
	minStorageSize    = 10
	maxServerActions  = 10
	PositionalArgHelp = "<UUID/Title/Hostname...>"
)

var CachedServers []upcloud.Server

func ServerCommand() commands.Command {
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

func SearchSingleServer(term string, service service.Server) (*upcloud.Server, error) {
	servers, err := SearchServer(&CachedServers, service, term, true)
	if err != nil {
		return nil, err
	}
	return servers[0], nil
}

func SearchServer(serversPtr *[]upcloud.Server, service service.Server, uuidOrHostnameOrTitle string, unique bool) ([]*upcloud.Server, error) {
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

func SearchAllServers(terms []string, service service.Server, unique bool) ([]string, error) {
	return commands.SearchResources(
		terms,
		func(term string) (interface{}, error) {
			return SearchServer(&CachedServers, service, term, unique)
		},
		func(in interface{}) string { return in.(*upcloud.Server).UUID })
}

func WaitForServerState(service service.Server, uuid, desiredState string, timeout time.Duration) (*upcloud.ServerDetails, error) {
	timer := time.After(timeout)
	for {
		time.Sleep(1 * time.Second)
		details, err := service.GetServerDetails(&request.GetServerDetailsRequest{UUID: uuid})
		if err != nil {
			return nil, err
		}
		switch details.State {
		case upcloud.ServerStateError:
			return nil, errors.New("server in error state")
		case desiredState:
			return details, nil
		}
		select {
		case <-timer:
			return nil, fmt.Errorf("timed out while waiting server to transition into %q", desiredState)
		default:
		}
	}
}

var WaitForServerFn = func(svc service.Server, state string, timeout time.Duration) func(uuid string, waitMsg string, err error) (interface{}, error) {
	return func(uuid string, waitMsg string, err error) (interface{}, error) {
		return WaitForServerState(svc, uuid, state, timeout)
	}
}

var getServerDetailsUuid = func(in interface{}) string { return in.(*upcloud.ServerDetails).UUID }

var getServerDetailsIpAddresses = func(in interface{}) []string {
	addresses := []string{}
	for _, iface := range in.(*upcloud.ServerDetails).Networking.Interfaces {
		for _, addr := range iface.IPAddresses {
			addresses = append(addresses, addr.Address)
		}
	}

	return addresses
}

type Request struct {
	ExactlyOne   bool
	BuildRequest func(server string) interface{}
	Service      service.Server
	Handler      ui.Handler
}

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

func GetArgCompFn(s service.Server) func(toComplete string) ([]string, cobra.ShellCompDirective) {
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
