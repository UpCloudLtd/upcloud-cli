package ip_address

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

const maxIpAddressActions = 10

func IpAddressCommand() commands.Command {
	return &ipAddressCommand{commands.New("ip-address", "Manage ip address")}
}

type ipAddressCommand struct {
	*commands.BaseCommand
}

func searchIpAddress(prtOrAddress string, service service.IpAddress) (*upcloud.IPAddress, error) {
	var result []*upcloud.IPAddress
	ips, err := service.GetIPAddresses()
	if err != nil {
		return nil, err
	}
	for _, ip := range ips.IPAddresses {
		if ip.Address == prtOrAddress || ip.PTRRecord == prtOrAddress {
			result = append(result, &ip)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no ip address was found with %s", prtOrAddress)
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("multiple ip addresses matched to query %q", prtOrAddress)
	}
	return result[0], nil
}

func searchIpAddresses(prtOrAddresses []string, service service.IpAddress) ([]*upcloud.IPAddress, error) {
	var result []*upcloud.IPAddress
	for _, prtOrAddress := range prtOrAddresses {
		ip, err := searchIpAddress(prtOrAddress, service)
		if err != nil {
			return nil, err
		}
		result = append(result, ip)
	}
	return result, nil
}

type Request struct {
	ExactlyOne   bool
	BuildRequest func(storage *upcloud.IPAddress) interface{}
	Service      service.IpAddress
	ui.HandleContext
}

func (s Request) Send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single ip address or ptr record is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one ip address or ptr record is required")
	}

	servers, err := searchIpAddresses(args, s.Service)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, server := range servers {
		requests = append(requests, s.BuildRequest(server))
	}

	return s.Handle(requests)
}
