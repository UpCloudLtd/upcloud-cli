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

type deleteCommand struct {
	*commands.BaseCommand
	service *service.Service
	index   int
}

func (s *deleteCommand) InitCommand() {
	s.SetPositionalArgHelp(server.PositionalArgHelp)
	s.ArgCompletion(server.GetArgCompFn(s.service))
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.index, "index", 0, "Interface index.")
	s.AddFlags(fs)
}

func DeleteCommand(service *service.Service) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a network interface"),
		service:     service,
	}
}

func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return server.Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				return &request.DeleteNetworkInterfaceRequest{
					ServerUUID: server.UUID,
					Index:      s.index,
				}
			},
			Service:    s.service,
			ExactlyOne: true,
			HandleContext: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.DeleteNetworkInterfaceRequest)
					return fmt.Sprintf("Deleting network interface %q of server %q", req.Index, req.ServerUUID)
				},
				MaxActions:    maxNetworkInterfaceActions,
				InteractiveUI: s.Config().InteractiveUI(),
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.service.DeleteNetworkInterface(req.(*request.DeleteNetworkInterfaceRequest))
				},
			},
		}.Send(args)
	}
}
