package networkinterface

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
	"strconv"
)

type createCommand struct {
	*commands.BaseCommand
	serverSvc  service.Server
	networkSvc service.Network
	params     createParams
}

// CreateCommand creates the "network-interface create" command
func CreateCommand(serverSvc service.Server, networkSvc service.Network) commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a network interface"),
		serverSvc:   serverSvc,
		networkSvc:  networkSvc,
	}
}

type createParams struct {
	req         request.CreateNetworkInterfaceRequest
	ipAddresses []string
	bootable    bool
	filtering   bool
	network     string
	family      string
}

var def = createParams{
	req: request.CreateNetworkInterfaceRequest{
		Type: upcloud.NetworkTypePrivate,
	},
	family: upcloud.IPAddressFamilyIPv4,
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.params.req = request.CreateNetworkInterfaceRequest{}
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.params.network, "network", def.network, "Virtual network ID or name to join.\n[Required]")
	fs.StringVar(&s.params.req.Type, "type", def.req.Type, "Set the type of the network.\nAvailable: public, utility, private")
	fs.StringVar(&s.params.family, "family", def.family, "The address family of new IP address.")
	fs.IntVar(&s.params.req.Index, "index", def.req.Index, "Interface index.")
	fs.BoolVar(&s.params.bootable, "bootable", def.bootable, "Whether to try booting through the interface.")
	fs.BoolVar(&s.params.filtering, "source-ip-filtering", def.filtering, "Whether source IP filtering is enabled on the interface. Disabling it is allowed only for SDN private interfaces.")
	fs.StringSliceVar(&s.params.ipAddresses, "ip-addresses", s.params.ipAddresses, "Array of IP addresses, multiple can be declared\n\n"+
		"Usage: --ip-addresses 94.237.112.143,94.237.112.144")
	s.AddFlags(fs)
}

func (s *createCommand) buildRequest() (*request.CreateNetworkInterfaceRequest, error) {
	if s.params.network == "" {
		s.params.req.IPAddresses = request.CreateNetworkInterfaceIPAddressSlice{{Family: s.params.family}}
	} else {

		if len(s.params.ipAddresses) == 0 {
			return nil, fmt.Errorf("ip-address is required")
		}
		ipAddresses, err := handleIPAddress(s.params.ipAddresses)
		if err != nil {
			return nil, err
		}
		s.params.req.IPAddresses = ipAddresses

		nw, err := network.SearchUniqueNetwork(s.params.network, s.networkSvc)
		if err != nil {
			return nil, err
		}
		s.params.req.NetworkUUID = nw.UUID
	}

	s.params.req.Bootable = upcloud.FromBool(s.params.bootable)
	s.params.req.SourceIPFiltering = upcloud.FromBool(s.params.filtering)
	return &s.params.req, nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		req, err := s.buildRequest()
		if err != nil {
			return nil, err
		}

		return server.Request{
			BuildRequest: func(uuid string) interface{} {
				req.ServerUUID = uuid
				return req
			},
			Service:    s.serverSvc,
			ExactlyOne: true,
			Handler: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.CreateNetworkInterfaceRequest)
					return fmt.Sprintf("Creating network interface for server %s network %s", req.ServerUUID, req.NetworkUUID)
				},
				ResultUUID:    func(in interface{}) string { return strconv.Itoa(in.(*upcloud.Interface).Index) },
				ResultPrefix:  "Index",
				MaxActions:    maxNetworkInterfaceActions,
				InteractiveUI: s.Config().InteractiveUI(),
				Action: func(req interface{}) (interface{}, error) {
					return s.networkSvc.CreateNetworkInterface(req.(*request.CreateNetworkInterfaceRequest))
				},
			},
		}.Send(args)
	}
}
