package router

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

var cachedRouters []upcloud.Router

const maxRouterActions = 10

// TODO: re-add
// const positionalArgHelp = "<UUID/Name...>"

// BaseRouterCommand creates the base "router" command
func BaseRouterCommand() commands.Command {
	return &routerCommand{commands.New("router", "Manage router")}
}

type routerCommand struct {
	*commands.BaseCommand
}

func searchRouter(term string, service service.Network, unique bool) ([]*upcloud.Router, error) {
	var result []*upcloud.Router

	if len(cachedRouters) == 0 {
		routers, err := service.GetRouters()
		if err != nil {
			return nil, err
		}
		cachedRouters = routers.Routers
	}

	for _, r := range cachedRouters {
		router := r
		if router.UUID == term || router.Name == term {
			result = append(result, &router)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no router was found with %s", term)
	}
	if len(result) > 1 && unique {
		return nil, fmt.Errorf("multiple routers matched to query %q", term)
	}
	return result, nil
}

func searchAllRouters(terms []string, service service.Network) ([]string, error) {
	return commands.SearchResources(
		terms,
		func(id string) (interface{}, error) {
			return searchRouter(id, service, true)
		},
		func(in interface{}) string { return in.(*upcloud.Router).UUID })
}

type routerRequest struct {
	ExactlyOne   bool
	BuildRequest func(uuid string) interface{}
	Service      service.Network
	Handler      ui.Handler
}

func (s routerRequest) send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single router uuid is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one router uuid is required")
	}

	servers, err := searchAllRouters(args, s.Service)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, server := range servers {
		requests = append(requests, s.BuildRequest(server))
	}

	return s.Handler.Handle(requests)
}
