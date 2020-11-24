package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

const maxNetworkActions = 10

func NetworkCommand() commands.Command {
	return &networkCommand{commands.New("network", "Manage network")}
}

type networkCommand struct {
	*commands.BaseCommand
}

var getNetworkUuid = func(in interface{}) string { return in.(*upcloud.Network).UUID }

func searchNetwork(uuidOrName string, service service.Network) (*upcloud.Network, error) {
	var result []*upcloud.Network
	networks, err := service.GetNetworks()
	if err != nil {
		return nil, err
	}
	for _, network := range networks.Networks {
		if network.UUID == uuidOrName || network.Name == uuidOrName {
			result = append(result, &network)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no network was found with %s", uuidOrName)
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("multiple networks matched to query %q", uuidOrName)
	}
	return result[0], nil
}

func searchNetworkes(uuidOrNames []string, service service.Network) ([]*upcloud.Network, error) {
	var result []*upcloud.Network
	for _, uuidOrName := range uuidOrNames {
		ip, err := searchNetwork(uuidOrName, service)
		if err != nil {
			return nil, err
		}
		result = append(result, ip)
	}
	return result, nil
}

type Request struct {
	ExactlyOne   bool
	BuildRequest func(storage *upcloud.Network) interface{}
	Service      service.Network
	ui.HandleContext
}

func (s Request) Send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single network uuid or name is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one network uuid or name is required")
	}

	servers, err := searchNetworkes(args, s.Service)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, server := range servers {
		requests = append(requests, s.BuildRequest(server))
	}

	return s.Handle(requests)
}
