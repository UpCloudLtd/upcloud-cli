package serverfirewall

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

type createCommand struct {
	*commands.BaseCommand
	serverSvc   service.Server
	firewallSvc service.Firewall
	params      createParams
}


func CreateCommand(serverSvc service.Server, firewallSvc service.Firewall) commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create firewall rule for server"),
		serverSvc:   serverSvc,
		firewallSvc: firewallSvc,
	}
}

/*var defaultCreateParams = &createParams{
	CreateFirewallRuleRequest: request.CreateFirewallRuleRequest{
	},
}
*/

var defaultCreateParams = request.CreateFirewallRuleRequest{
}

type createParams struct {
	req                        request.CreateFirewallRuleRequest
	direction                  string
	action                     string
	position                   string
	family                     string
	protocol                   string
	icmp_type                  string
	destination_address_start  string
	destination_address_end    string
	destination_port_start     string
	destination_port_end       string
	source_address_start       string
	source_address_end         string
	source_port_start          string
	source_port_end            string
	comment                    string

}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	// s.params = createParams{firewallSvc.CreateFirewallRuleRequest: request.CreateFirewallRuleRequest{}}

	def := defaultCreateParams

	flagSet.StringVar(&s.params.direction, "direction", def.FirewallRule.Direction, "")
	flagSet.StringVar(&s.params.action, "action", def.FirewallRule.Action, "")
	flagSet.StringVar(&s.params.family, "family", def.FirewallRule.Family, "")
/*	flagSet.StringVar(&s.params.Position, "position", def.Position, "")
	flagSet.StringVar(&s.params.Protocol, "protocol", def.Protocol, "")
	flagSet.StringVar(&s.params.Icmp_type, "icmp_type", def.Icmp_type, "")
	flagSet.StringVar(&s.params.Destination_address_start, "destination_address_start", def.Destination_address_start, "")
	flagSet.StringVar(&s.params.Destination_address_end, "destination_address_end", def.Destination_address_end, "")
	flagSet.StringVar(&s.params.Destination_port_start, "destination_port_start", def.Destination_port_start, "")
	flagSet.StringVar(&s.params.Destination_port_end, "destination_port_end", def.Destination_port_end, "")
	flagSet.StringVar(&s.params.Source_address_start, "source_address_start", def.Source_address_start, "")
	flagSet.StringVar(&s.params.Source_address_end, "source_address_end", def.Source_address_end, "")
	flagSet.StringVar(&s.params.Source_port_start, "source_port_start", def.Source_port_start, "")
	flagSet.StringVar(&s.params.Source_address_end, "source_address_end", def.Source_address_end, "")
	flagSet.StringVar(&s.params.Comment, "comment", def.Comment, "")*/

	s.AddFlags(flagSet)
}
	
// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		s.params.direction = upcloud.FirewallRuleDirectionIn

/*		if s.params.direction == "" {
			return nil, fmt.Errorf("direction is required")
		}

		if s.params.action == "" {
			return nil, fmt.Errorf("action is required")
		}

		if s.params.family == "" {
			return nil, fmt.Errorf("family is required")
		}

		if s.params.destination_address_start == "" && s.params.destination_address_end != "" {
			return nil, fmt.Errorf("destination_address_start is required if destination_address_end is set")
		}

		if s.params.destination_address_end == "" && s.params.destination_address_start != "" {
			return nil, fmt.Errorf("destination_address_end is required if destination_address_start is set")
		}

		if s.params.destination_port_start == "" && s.params.destination_port_end != "" {
			return nil, fmt.Errorf("destination_port_start is required if destination_port_end is set")
		}

		if s.params.destination_port_end == "" && s.params.destination_port_start != "" {
			return nil, fmt.Errorf("destination_port_end is required if destination_port_start is set")
		}

		if s.params.source_address_start == "" && s.params.source_address_end != "" {
			return nil, fmt.Errorf("source_address_start is required if source_address_end is set")
		}

		if s.params.source_address_end == "" && s.params.source_address_start != "" {
			return nil, fmt.Errorf("source_address_end is required if source_address_start is set")
		}

		if s.params.source_port_start == "" && s.params.source_port_end != "" {
			return nil, fmt.Errorf("source_port_start is required if source_port_end is set")
		}

		if s.params.source_port_end == "" && s.params.source_port_start != "" {
			return nil, fmt.Errorf("source_port_end is required if source_port_start is set")
		}
*/
/*		return ui.HandleContext{
			RequestID:       func(in interface{}) string { return in.(*request.CreateFirewallRuleRequest).Hostname },
			ResultUUID:      getServerDetailsUUID,
			ResultExtras:    getServerDetailsIPAddresses,
			ResultExtraName: "IP addresses",
			InteractiveUI:   s.Config().InteractiveUI(),
			WaitMsg:         "server starting",
			WaitFn:          waitForServer(s.serverSvc, upcloud.ServerStateStarted, s.Config().ClientTimeout()),
			MaxActions:      5,
			ActionMsg:       "Creating server",
			Action: func(req interface{}) (interface{}, error) {
				return s.serverSvc.CreateServer(req.(*request.CreateFirewallRuleRequest))
			},
		}.Handle(commands.ToArray(&req))*/

		return server.Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.req
				req.ServerUUID = uuid
				return &req
			},
			Service:    s.serverSvc,
			ExactlyOne: true,
			Handler: ui.HandleContext{
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    10,
				MessageFn: func(in interface{}) string {
					req := in.(*request.CreateFirewallRuleRequest)
					return fmt.Sprintf("Creating firewall rule for server %q", req.ServerUUID)
				},
				Action: func(req interface{}) (interface{}, error) {
					return s.firewallSvc.CreateFirewallRule(req.(*request.CreateFirewallRuleRequest))
				},
			},
		}.Send(args)
	}
}
