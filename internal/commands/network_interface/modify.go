package network_interface

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	service     *service.Service
	req         request.ModifyNetworkInterfaceRequest
	bootable    string
	filtering   string
	ipAddresses []string
}

func ModifyCommand(service *service.Service) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a network interface"),
		service:     service,
	}
}

func (s *modifyCommand) BuildRequest() (*request.ModifyNetworkInterfaceRequest, error) {
	ipAddresses, err := handleIpAddress(s.ipAddresses)
	if err != nil {
		return nil, err
	}
	s.req.IPAddresses = ipAddresses

	if s.bootable != "" {
		s.req.Bootable = upcloud.FromBool(s.bootable == "true")
	}
	if s.filtering != "" {
		s.req.SourceIPFiltering = upcloud.FromBool(s.filtering == "true")
	}
	return &s.req, nil
}

func (s *modifyCommand) InitCommand() {
	s.SetPositionalArgHelp(server.PositionalArgHelp)
	s.ArgCompletion(server.GetArgCompFn(s.service))
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.req.CurrentIndex, "index", 0, "Index of the interface to modify.")
	fs.IntVar(&s.req.NewIndex, "new-index", 0, "Index of the interface to modify.")
	fs.StringVar(&s.bootable, "bootable", "", "Whether to try booting through the interface.")
	fs.StringVar(&s.filtering, "source-ip-filtering", "", "Whether source IP filtering is enabled on the interface. Disabling it is allowed only for SDN private interfaces.")
	fs.StringArrayVar(&s.ipAddresses, "ip-address", s.ipAddresses, "Array of IP addresses")
	s.AddFlags(fs)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return server.Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				return &request.ModifyNetworkInterfaceRequest{ServerUUID: server.UUID}
			},
			Service:    s.service,
			ExactlyOne: true,
			HandleContext: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.ModifyNetworkInterfaceRequest)
					return fmt.Sprintf("Modifying network interface %q of server %q", req.CurrentIndex, req.ServerUUID)
				},
				MaxActions:    maxNetworkInterfaceActions,
				InteractiveUI: s.Config().InteractiveUI(),
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyNetworkInterface(req.(*request.ModifyNetworkInterfaceRequest))
				},
			},
		}.Send(args)
	}
}
