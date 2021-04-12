package serverfirewall

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"net"
	"strings"
	"sync"

	"github.com/UpCloudLtd/cli/internal/commands"
)

// ShowCommand is the 'server firewall show' command, displaying firewall details
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show server firewall details."),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	// TODO: reimplement
	//	s.SetPositionalArgHelp(positionalArgHelp)
}

type fwRuleAddress struct {
	AddressStart string `json:"address_start,omitempty"`
	AddressEnd   string `json:"address_end,omitempty"`
	PortStart    string `json:"port_start,omitempty"`
	PortEnd      string `json:"port_end,omitempty"`
}

type fwProto struct {
	Family   string `json:"family,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	ICMPType string `json:"icmptype,omitempty"`
}

// Execute implements Command.Execute
func (s *showCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// TODO(aakso): implement prompting with readline support
	if arg == "" {
		return nil, fmt.Errorf("one server hostname, title or uuid is required")
	}
	// get rules and server details in parallel
	var wg sync.WaitGroup
	var rules *upcloud.FirewallRules
	var fwRuleErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		rules, fwRuleErr = exec.Firewall().GetFirewallRules(&request.GetFirewallRulesRequest{ServerUUID: arg})
	}()
	server, err := exec.Server().GetServerDetails(&request.GetServerDetailsRequest{UUID: arg})
	if err != nil {
		return nil, err
	}
	wg.Wait()
	if fwRuleErr != nil {
		return nil, err
	}
	// build output
	var fwRows []output.TableRow
	for _, rule := range rules.FirewallRules {
		fwRows = append(fwRows, output.TableRow{
			rule.Position,
			rule.Action,
			fwRuleAddress{
				rule.SourceAddressStart,
				rule.SourceAddressEnd,
				rule.SourcePortStart,
				rule.SourcePortEnd,
			},
			fwRuleAddress{
				rule.DestinationAddressStart,
				rule.DestinationAddressEnd,
				rule.DestinationPortStart,
				rule.DestinationPortEnd,
			},
			rule.Direction,
			fwProto{
				rule.Family,
				rule.Protocol,
				rule.ICMPType,
			},
		})
	}
	return output.Combined{
		output.CombinedSection{
			Key:   "rules",
			Title: "Firewall rules",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "index", Header: "#"},
					{Key: "action", Header: "Action", Format: actionFormat},
					{Key: "source", Header: "Source", Format: addressFormat},
					{Key: "destination", Header: "Destination", Format: addressFormat},
					{Key: "direction", Header: "Dir"},
					{Key: "protocol", Header: "Proto", Format: protoFormat},
				},
				Rows: fwRows,
			},
		},
		output.CombinedSection{
			Contents: output.Details{Sections: []output.DetailSection{
				{Rows: []output.DetailRow{
					{Key: "enabled", Title: "Enabled", Value: server.Firewall == "on", Format: output.BoolFormat},
				}},
			}},
		},
	}, nil
}

func protoFormat(val interface{}) (text.Colors, string, error) {
	if fwp, ok := val.(fwProto); ok {
		return nil, formatProto(fwp), nil
	}
	return nil, fmt.Sprint(val), nil
}

func addressFormat(val interface{}) (text.Colors, string, error) {
	if fwa, ok := val.(fwRuleAddress); ok {
		return ui.DefaultAddressColours, formatMatch(fwa), nil
	}
	return ui.DefaultAddressColours, fmt.Sprint(val), nil
}

func actionFormat(val interface{}) (text.Colors, string, error) {
	if actStr, ok := val.(string); ok {
		if actStr == "drop" {
			return text.Colors{text.FgHiRed}, actStr, nil
		}
		return text.Colors{text.FgHiGreen}, actStr, nil
	}
	return nil, fmt.Sprint(val), nil
}

func formatMatch(address fwRuleAddress) string {
	var sb strings.Builder
	ipStart := net.ParseIP(address.AddressStart)
	ipStop := net.ParseIP(address.AddressEnd)
	if ipStart != nil {
		if ipStart.Equal(ipStop) {
			// TODO: ermm, reimplement.. when we figure out if this is really needed + how to do it
			// sb.WriteString(ui.DefaultAddressColours.Sprint(ipStart))
			sb.WriteString(fmt.Sprint(ipStart))
		} else {
			// sb.WriteString(ui.DefaultAddressColours.Sprintf("%s →\n%s", ipStart, ipStop))
			sb.WriteString(fmt.Sprintf("%s →\n%s", ipStart, ipStop))
		}
	}
	if address.PortStart != "" {
		if ipStart != nil {
			sb.WriteString("\n")
		}
		if address.PortStart == address.PortEnd {
			sb.WriteString(fmt.Sprintf("port: %s", address.PortStart))
		} else {
			sb.WriteString(fmt.Sprintf("port: %s → %s", address.PortStart, address.PortEnd))
		}
	}
	return sb.String()
}

func formatProto(proto fwProto) string {
	if proto.Protocol == "" {
		return proto.Family
	}
	if proto.ICMPType == "" {
		return fmt.Sprintf("%s/%s", proto.Family, proto.Protocol)
	}
	return fmt.Sprintf("%s/%s/%s", proto.Family, proto.Protocol, proto.ICMPType)
}
