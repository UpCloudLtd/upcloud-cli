package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/validation"
)

const (
	minStorageSize   = 10
	maxServerActions = 10
)

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

func searchServer(serversPtr *[]upcloud.Server, service *service.Service, uuidOrHostnameOrTitle string, unique bool) (*upcloud.Server, error) {
	if serversPtr == nil || service == nil {
		return nil, fmt.Errorf("no servers or service passed")
	}
	servers := *serversPtr
	if err := validation.Uuid4(uuidOrHostnameOrTitle); err != nil || servers == nil {
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
	return matched[0], nil
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

func WaitForServerState(service *service.Service, uuid, desiredState string, timeout time.Duration) (*upcloud.ServerDetails, error) {
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
