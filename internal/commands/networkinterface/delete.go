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

type deleteCommand struct {
	*commands.BaseCommand
	networkSvc service.Network
	serverSvc  service.Server
	index      int
}

// DeleteCommand creates the "network-interface delete" command
func DeleteCommand(networkSvc service.Network, serverSvc service.Server) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a network interface"),
		networkSvc:  networkSvc,
		serverSvc:   serverSvc,
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	s.SetPositionalArgHelp(server.PositionalArgHelp)
	s.ArgCompletion(server.GetServerArgumentCompletionFunction(s.serverSvc))
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.index, "index", 0, "Interface index.")
	s.AddFlags(fs)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.index == 0 {
			return nil, fmt.Errorf("index is required")
		}

		return server.Request{
			BuildRequest: func(uuid string) interface{} {
				return &request.DeleteNetworkInterfaceRequest{
					ServerUUID: uuid,
					Index:      s.index,
				}
			},
			Service:    s.serverSvc,
			ExactlyOne: true,
			Handler: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.DeleteNetworkInterfaceRequest)
					return fmt.Sprintf("Deleting network interface %d of server %q", req.Index, req.ServerUUID)
				},
				MaxActions:    maxNetworkInterfaceActions,
				InteractiveUI: s.Config().InteractiveUI(),
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.networkSvc.DeleteNetworkInterface(req.(*request.DeleteNetworkInterfaceRequest))
				},
			},
		}.Send(args)
	}
}
