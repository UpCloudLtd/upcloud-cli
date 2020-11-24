package router

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

const maxRouterActions = 10

func RouterCommand() commands.Command {
	return &routerCommand{commands.New("router", "Manage router")}
}

type routerCommand struct {
	*commands.BaseCommand
}

var getRouterUuid = func(in interface{}) string { return in.(*upcloud.Router).UUID }

func searchRouter(uuidOrName string, service *service.Service) (*upcloud.Router, error) {
	var result []*upcloud.Router
	routers, err := service.GetRouters()
	if err != nil {
		return nil, err
	}
	for _, router := range routers.Routers {
		if router.UUID == uuidOrName || router.Name == uuidOrName {
			result = append(result, &router)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no router was found with %s", uuidOrName)
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("multiple routers matched to query %q", uuidOrName)
	}
	return result[0], nil
}

func searchRouteres(uuidOrNames []string, service *service.Service) ([]*upcloud.Router, error) {
	var result []*upcloud.Router
	for _, uuidOrName := range uuidOrNames {
		ip, err := searchRouter(uuidOrName, service)
		if err != nil {
			return nil, err
		}
		result = append(result, ip)
	}
	return result, nil
}

type Request struct {
	ExactlyOne   bool
	BuildRequest func(storage *upcloud.Router) interface{}
	Service      *service.Service
	ui.HandleContext
}

func (s Request) Send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single router uuid is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one router uuid is required")
	}

	servers, err := searchRouteres(args, s.Service)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, server := range servers {
		requests = append(requests, s.BuildRequest(server))
	}

	return s.Handle(requests)
}
