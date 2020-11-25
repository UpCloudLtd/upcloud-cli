package network_interface

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/network"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	service *service.Service
	params  createParams
}

func CreateCommand(service *service.Service) commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a network"),
		service:     service,
	}
}

type createParams struct {
	req         request.CreateNetworkInterfaceRequest
	ipAddresses []string
	bootable    bool
	filtering   bool
	network     string
}

var def = createParams{
	req: request.CreateNetworkInterfaceRequest{},
}

func (s *createCommand) InitCommand() {
	s.params.req = request.CreateNetworkInterfaceRequest{}
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.params.network, "network", def.network, "Virtual network ID to join.")
	fs.StringVar(&s.params.req.Type, "type", def.req.Type, "Set the type of the network.\nAvailable: public, utility, private")
	fs.IntVar(&s.params.req.Index, "index", def.req.Index, "Interface index.")
	fs.BoolVar(&s.params.bootable, "bootable", def.bootable, "Whether to try booting through the interface.")
	fs.BoolVar(&s.params.filtering, "source-ip-filtering", def.filtering, "Whether source IP filtering is enabled on the interface. Disabling it is allowed only for SDN private interfaces.")
	fs.StringArrayVar(&s.params.ipAddresses, "ip-address", s.params.ipAddresses, "Array of IP addresses")
	s.AddFlags(fs)
}

func (s *createCommand) BuildRequest() (*request.CreateNetworkInterfaceRequest, error) {
	ipAddresses, err := handleIpAddress(s.params.ipAddresses)
	if err != nil {
		return nil, err
	}
	nw, err := network.SearchNetwork(s.params.network, s.service)
	if err != nil {
		return nil, err
	}
	s.params.req.NetworkUUID = nw.UUID
	s.params.req.IPAddresses = ipAddresses
	s.params.req.Bootable = upcloud.FromBool(s.params.bootable)
	s.params.req.SourceIPFiltering = upcloud.FromBool(s.params.filtering)
	return &s.params.req, nil
}

func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		req, err := s.BuildRequest()
		if err != nil {
			return nil, err
		}

		return server.Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req.ServerUUID = server.UUID
				return req
			},
			Service:    s.service,
			ExactlyOne: true,
			HandleContext: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.CreateNetworkInterfaceRequest)
					return fmt.Sprintf("Creating network interface for server %s network %s", req.ServerUUID, req.NetworkUUID)
				},
				ResultUUID:    func(in interface{}) string { return fmt.Sprintf("Index %d", in.(*upcloud.Interface).Index) },
				MaxActions:    maxNetworkInterfaceActions,
				InteractiveUI: s.Config().InteractiveUI(),
				Action: func(req interface{}) (interface{}, error) {
					return s.service.CreateNetworkInterface(req.(*request.CreateNetworkInterfaceRequest))
				},
			},
		}.Send(args)
	}
}
