package server

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/ui"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/cli/internal/commands"
)

const (
	minStorageSize   = 10
	maxServerActions = 10
)

var cachedServers []upcloud.Server

func ServerCommand() commands.Command {
	return &serverCommand{
		Command: commands.New("server", "List, show & control servers"),
	}
}

type serverCommand struct {
	commands.Command
}

type ActionResult struct {
	Error error
	Uuid  string
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

func SearchServer(serversPtr *[]upcloud.Server, service service.Server, uuidOrHostnameOrTitle string, unique bool) ([]*upcloud.Server, error) {
	if serversPtr == nil || service == nil {
		return nil, fmt.Errorf("no servers or service passed")
	}
	servers := *serversPtr
	if len(cachedServers) == 0 {
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

func SearchAllArgs(uuidOrTitle []string, service service.Server, unique bool) ([]*upcloud.Server, error) {
	var result []*upcloud.Server
	for _, id := range uuidOrTitle {
		matchedResults, err := SearchServer(&cachedServers, service, id, unique)
		if err != nil {
			return nil, err
		}
		result = append(result, matchedResults...)
	}
	return result, nil
}

func StateColour(state string) text.Colors {
	switch state {
	case upcloud.ServerStateStarted:
		return text.Colors{text.FgGreen}
	case upcloud.ServerStateError:
		return text.Colors{text.FgHiRed, text.Bold}
	case upcloud.ServerStateMaintenance:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

func WaitForServerState(service service.Server, uuid, desiredState string, timeout time.Duration) (*upcloud.ServerDetails, error) {
	timer := time.After(timeout)
	for {
		time.Sleep(5 * time.Second)
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

type ServerFirewall interface {
	service.Server
	service.Firewall
}

type Request struct {
	ExactlyOne   bool
	BuildRequest func(storage *upcloud.Server) interface{}
	Service      service.Server
	ui.HandleContext
}

func (s Request) Send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single server uuid is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one server uuid is required")
	}

	servers, err := SearchAllArgs(args, s.Service, true)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, server := range servers {
		requests = append(requests, s.BuildRequest(server))
	}

	return s.Handle(requests)
}
