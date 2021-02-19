package networkinterface

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	networkSvc  service.Network
	serverSvc   service.Server
	req         request.ModifyNetworkInterfaceRequest
	bootable    string
	filtering   string
	ipAddresses []string
}

// ModifyCommand creates the "network-interface modify" command
func ModifyCommand(networkSvc service.Network, serverSvc service.Server) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a network interface"),
		serverSvc:   serverSvc,
		networkSvc:  networkSvc,
	}
}

func (s *modifyCommand) buildRequest() (*request.ModifyNetworkInterfaceRequest, error) {
	ipAddresses, err := handleIPAddress(s.ipAddresses)
	if err != nil {
		return nil, err
	}
	s.req.IPAddresses = ipAddresses

	if s.bootable != "" {
		bootable, err := commands.BoolFromString(s.bootable)
		if err != nil {
			return nil, err
		}
		s.req.Bootable = *bootable
	}
	if s.filtering != "" {
		filtering, err := commands.BoolFromString(s.filtering)
		if err != nil {
			return nil, err
		}
		s.req.SourceIPFiltering = *filtering
	}
	return &s.req, nil
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	s.SetPositionalArgHelp(server.PositionalArgHelp)
	s.ArgCompletion(server.GetServerArgumentCompletionFunction(s.serverSvc))
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.req.CurrentIndex, "index", s.req.CurrentIndex, "Index of the interface to modify. [Required]")
	fs.IntVar(&s.req.NewIndex, "new-index", s.req.NewIndex, "Index of the interface to modify.")
	fs.StringVar(&s.bootable, "bootable", s.bootable, "Whether to try booting through the interface.")
	fs.StringVar(&s.filtering, "source-ip-filtering", s.filtering, "Whether source IP filtering is enabled on the interface. Disabling it is allowed only for SDN private interfaces.")
	fs.StringSliceVar(&s.ipAddresses, "ip-addresses", s.ipAddresses, "Array of IP addresses, multiple can be declared\nUsage: --ip-address address=94.237.112.143,family=IPv4")
	s.AddFlags(fs)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if s.req.CurrentIndex == 0 {
			return nil, fmt.Errorf("index is required")
		}

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
					req := in.(*request.ModifyNetworkInterfaceRequest)
					return fmt.Sprintf("Modifying network interface %q of server %q", req.CurrentIndex, req.ServerUUID)
				},
				MaxActions:    maxNetworkInterfaceActions,
				InteractiveUI: s.Config().InteractiveUI(),
				Action: func(req interface{}) (interface{}, error) {
					return s.networkSvc.ModifyNetworkInterface(req.(*request.ModifyNetworkInterfaceRequest))
				},
			},
		}.Send(args)
	}
}
