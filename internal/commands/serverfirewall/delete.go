package serverfirewall

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
	serverSvc  service.Server
	firewallSvc service.Firewall
	params     deleteParams
}

func DeleteCommand(serverSvc service.Server, firewallSvc service.Firewall) commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Removes a firewall rule from a server. Firewall rules must be removed individually. The positions of remaining firewall rules will be adjusted after a rule is removed."),
		serverSvc:   serverSvc,
		firewallSvc:  firewallSvc,
	}
}

var defaultRemoveParams = request.DeleteFirewallRuleRequest{
}

type deleteParams struct {
	request.DeleteFirewallRuleRequest
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	def := defaultRemoveParams

	flagSet.IntVar(&s.params.Position, "position", def.Position, "1-1000")

	s.AddFlags(flagSet)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *deleteCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Position == 0 {
			return nil, fmt.Errorf("Position is required.")
		}

		return server.Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.DeleteFirewallRuleRequest
				req.ServerUUID = uuid
				return &req
			},
			Service: s.serverSvc,
			Handler: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.DeleteFirewallRuleRequest)
					return fmt.Sprintf("Remove firewall rule at position %d from server %q", s.params.Position, req.ServerUUID)
				},
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    10,
				Action: func(req interface{}) (interface{}, error) {
					return nil, s.firewallSvc.DeleteFirewallRule(req.(*request.DeleteFirewallRuleRequest))
				},
			},
		}.Send(args)
	}
}

