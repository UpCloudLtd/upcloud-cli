package serverfirewall
/*
import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type removeCommand struct {
	*commands.BaseCommand
	serverSvc  service.Server
	firewallSvc service.Firewall
	params     removeParams
}


var defaultRemoveParams = &RemoveParams{
	RemoveFirewallRuleRequest: request.RemoveFirewallRuleRequest{},
}

type removeParams struct {
	req    request.RemoveFirewallRuleRequest
}

func RemoveCommand(serverSvc service.Server, firewallSvc service.Firewall) commands.Command {
	return &detachCommand{
		BaseCommand: commands.New("remove", "Remove firewall"),
		serverSvc:   serverSvc,
		firewallSvc:  firewallSvc,
	}
}

// InitCommand implements Command.InitCommand
func (s *detachCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(server.GetServerArgumentCompletionFunction(s.serverSvc))
	s.params = detachParams{RemoveFirewallRuleRequest: request.RemoveFirewallRuleRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Address, "address", defaultDetachParams.Address, "Detach the storage attached to this address.\n[Required]")

	s.AddFlags(flagSet)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *detachCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Address == "" {
			return nil, fmt.Errorf("address is required")
		}

		return server.Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.RemoveFirewallRuleRequest
				req.ServerUUID = uuid
				return &req
			},
			Service: s.serverSvc,
			Handler: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.RemoveFirewallRuleRequest)
					return fmt.Sprintf("Detaching address %q from server %q", req.Address, req.ServerUUID)
				},
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				Action: func(req interface{}) (interface{}, error) {
					return s.firewallSvc.RemoveFirewallRule(req.(*request.RemoveFirewallRuleRequest))
				},
			},
		}.Send(args)
	}
}
*/
