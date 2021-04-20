package serverfirewall

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/m7shapan/cidr"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	direction            string
	family               string
	action               string
	position             int
	protocol             string
	icmpType             string
	destinationIPBlock   string
	destinationPortStart string
	destinationPortEnd   string
	sourceIPBlock        string
	sourcePortStart      string
	sourcePortEnd        string
	comment              string
	completion.Server
	resolver.CachingServer
}

// CreateCommand creates the "server filewall create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a new firewall rule", ""),
	}
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return 10
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	flagSet.StringVar(&s.direction, "direction", "", "Rule direction. Available: in / out")
	flagSet.StringVar(&s.action, "action", "", "Rule action. Available: accept / drop")
	flagSet.StringVar(&s.family, "family", "", "IP family. Available: IPv4, IPv6")
	flagSet.IntVar(&s.position, "position", 0, "Position in relation to other rules. Available: 1-1000")
	flagSet.StringVar(&s.protocol, "protocol", "", "Protocol. Available: tcp, udp, icmp")
	flagSet.StringVar(&s.icmpType, "icmp-type", "", "ICMP type. Available: 0-255")
	flagSet.StringVar(&s.destinationIPBlock, "dest-ipaddress-block", "", "Destination IP address block.")
	flagSet.StringVar(&s.destinationPortStart, "destination-port-start", "", "Destination port range start. Available: 1-65535")
	flagSet.StringVar(&s.destinationPortEnd, "destination-port-end", "", "Destination port range end.")
	flagSet.StringVar(&s.sourceIPBlock, "src-ipaddress-block", "", "Source IP address block.")
	flagSet.StringVar(&s.sourcePortStart, "source-port-start", "", "Source port range start.")
	flagSet.StringVar(&s.sourcePortEnd, "source-port-end", "", "Destination port range end.")
	flagSet.StringVar(&s.comment, "comment", "", "Freeform comment that can include 0-250 characters.")

	s.AddFlags(flagSet)
}

// Execute implements commands.MultipleArgumentCommand
func (s *createCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if s.direction == "" {
		return nil, fmt.Errorf("direction is required")
	}

	if s.action == "" {
		return nil, fmt.Errorf("action is required")
	}

	if s.family == "" {
		return nil, fmt.Errorf("family (IPv4/IPv6) is required")
	}

	if s.family != "IPv4" && s.family != "IPv6" {
		return nil, fmt.Errorf("invalid family, use either IPv4 or IPv6")
	}

	if s.destinationPortStart == "" && s.destinationPortEnd != "" {
		return nil, fmt.Errorf("destination-port-start is required if destination-port-end is set")
	}

	if s.destinationPortEnd == "" && s.destinationPortStart != "" {
		return nil, fmt.Errorf("destination-port-end is required if destination-port-start is set")
	}

	if s.sourcePortStart == "" && s.sourcePortEnd != "" {
		return nil, fmt.Errorf("source-port-start is required if source-port-end is set")
	}

	if s.sourcePortEnd == "" && s.sourcePortStart != "" {
		return nil, fmt.Errorf("source-port-end is required if source-port-start is set")
	}

	var (
		destinationNetwork      *cidr.ParsedCIDR
		sourceNetwork           *cidr.ParsedCIDR
		destinationAddressStart string
		destinationAddressEnd   string
		sourceAddressStart      string
		sourceAddressEnd        string
		err                     error
	)
	if s.destinationIPBlock != "" {
		destinationNetwork, err = cidr.ParseCIDR(s.destinationIPBlock)
		if err != nil {
			return nil, fmt.Errorf("dest-ipaddress-block parse error: %s", err)
		}
		destinationAddressStart = destinationNetwork.FirstIP.String()
		destinationAddressEnd = destinationNetwork.LastIP.String()
	}
	if s.sourceIPBlock != "" {
		sourceNetwork, err = cidr.ParseCIDR(s.sourceIPBlock)
		if err != nil {
			return nil, fmt.Errorf("src-ipaddress-block parse error: %s", err)
		}
		sourceAddressStart = sourceNetwork.FirstIP.String()
		sourceAddressEnd = sourceNetwork.LastIP.String()
	}

	msg := fmt.Sprintf("creating firewall rule for server %v", arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	res, err := exec.Firewall().CreateFirewallRule(&request.CreateFirewallRuleRequest{
		ServerUUID: arg,
		FirewallRule: upcloud.FirewallRule{
			Action:                  s.action,
			Comment:                 s.comment,
			DestinationAddressStart: destinationAddressStart,
			DestinationAddressEnd:   destinationAddressEnd,
			DestinationPortStart:    s.destinationPortStart,
			DestinationPortEnd:      s.destinationPortEnd,
			Direction:               s.direction,
			Family:                  s.family,
			ICMPType:                s.icmpType,
			Position:                s.position,
			Protocol:                s.protocol,
			SourceAddressStart:      sourceAddressStart,
			SourceAddressEnd:        sourceAddressEnd,
			SourcePortStart:         s.sourcePortStart,
			SourcePortEnd:           s.sourcePortEnd,
		},
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil

}
