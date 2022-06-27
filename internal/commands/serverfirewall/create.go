package serverfirewall

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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
		BaseCommand: commands.New(
			"create",
			"Create a new firewall rule",
			"upctl server firewall create 00038afc-d526-4148-af0e-d2f1eeaded9b --direction in --action drop",
			"upctl server firewall create 00038afc-d526-4148-af0e-d2f1eeaded9b --direction in --action accept --family IPv4",
			"upctl server firewall create 00038afc-d526-4148-af0e-d2f1eeaded9b --direction in --action drop --family IPv4 --src-ipaddress-block 10.11.0.88/24",
		),
	}
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return 10
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.Cobra().Long = `Create a new firewall rule

To edit the default rule of the firewall, set only --direction and --action
parameters. This creates catch-all rule that will take effect when no other rule
matches. Note that the default rule must be positioned after all other rules.
Use --position parameter or create default rule after other rules.`

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
	s.Cobra().MarkFlagRequired("direction") //nolint:errcheck
	s.Cobra().MarkFlagRequired("action")    //nolint:errcheck
	s.Cobra().MarkFlagsRequiredTogether("destination-port-start", "destination-port-end")
	s.Cobra().MarkFlagsRequiredTogether("source-port-start", "source-port-end")
}

// Execute implements commands.MultipleArgumentCommand
func (s *createCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if s.family != "" && s.family != "IPv4" && s.family != "IPv6" {
		return nil, fmt.Errorf("invalid family, use either IPv4 or IPv6")
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
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}
	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
