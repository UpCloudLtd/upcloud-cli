package serverfirewall

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/m7shapan/cidr"
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
		BaseCommand: commands.New("create", "Creates a new firewall rule."),
		serverSvc:   serverSvc,
		firewallSvc: firewallSvc,
	}
}

var (
	defaultCreateParams = request.CreateFirewallRuleRequest{}
)

type createParams struct {
	request.CreateFirewallRuleRequest
	DestIPBlock string
	SrcIPBlock  string
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	def := defaultCreateParams

	flagSet.StringVar(&s.params.Direction, "direction", def.FirewallRule.Direction, "in / out")
	flagSet.StringVar(&s.params.Action, "action", def.FirewallRule.Action, "accept / drop")
	flagSet.StringVar(&s.params.Family, "family", def.FirewallRule.Family, "IPv4 / IPv6")
	flagSet.IntVar(&s.params.Position, "position", def.Position, "1-1000")
	flagSet.StringVar(&s.params.Protocol, "protocol", def.Protocol, "tcp / udp / icmp")
	flagSet.StringVar(&s.params.ICMPType, "icmp-type", def.ICMPType, "0-255")
	flagSet.StringVar(&s.params.DestIPBlock, "dest-ipaddress-block", "", "Valid IP address")
	flagSet.StringVar(&s.params.DestinationPortStart, "destination-port-start", def.DestinationPortStart, "1-65535")
	flagSet.StringVar(&s.params.DestinationPortEnd, "destination-port-end", def.DestinationPortEnd, "1-65535")
	flagSet.StringVar(&s.params.SrcIPBlock, "src-ipaddress-block", "", "Valid IP address")
	flagSet.StringVar(&s.params.SourcePortStart, "source-port-start", def.SourcePortStart, "1-65535")
	flagSet.StringVar(&s.params.SourcePortEnd, "source-port-end", def.SourcePortEnd, "1-65535")
	flagSet.StringVar(&s.params.Comment, "comment", def.Comment, "0-250 characters")

	s.AddFlags(flagSet)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *createCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Direction == "" {
			return nil, fmt.Errorf("Direction is required.")
		}

		if s.params.Action == "" {
			return nil, fmt.Errorf("Action is required.")
		}

		if s.params.Family == "" {
			return nil, fmt.Errorf("Family is required. Use 'IPv4' or 'IPv6'.")
		}

		if s.params.Family != "IPv4" && s.params.Family != "IPv6" {
			return nil, fmt.Errorf("Invalid Family. 'IPv4' / 'IPv6'.")
		}

		NetDst, err := cidr.ParseCIDR(s.params.DestIPBlock)
		if err != nil {
			return nil, fmt.Errorf("dest-ipaddress-block parse error: %s", err)
		}
		s.params.DestinationAddressStart = NetDst.FirstIP.String()
		s.params.DestinationAddressEnd = NetDst.LastIP.String()

		if s.params.DestinationPortStart == "" && s.params.DestinationPortEnd != "" {
			return nil, fmt.Errorf("destination-port-start is required if destination-port-end is set")
		}

		if s.params.DestinationPortEnd == "" && s.params.DestinationPortStart != "" {
			return nil, fmt.Errorf("destination-port-end is required if destination-port-start is set")
		}

		NetSrc, err := cidr.ParseCIDR(s.params.SrcIPBlock)
		if err != nil {
			return nil, fmt.Errorf("src-ipaddress-block parse error: %s", err)
		}
		s.params.SourceAddressStart = NetSrc.FirstIP.String()
		s.params.SourceAddressEnd = NetSrc.LastIP.String()

		if s.params.SourcePortStart == "" && s.params.SourcePortEnd != "" {
			return nil, fmt.Errorf("source-port-start is required if source-port-end is set")
		}

		if s.params.SourcePortEnd == "" && s.params.SourcePortStart != "" {
			return nil, fmt.Errorf("source-port-end is required if source-port-start is set")
		}

		return server.Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.CreateFirewallRuleRequest
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
