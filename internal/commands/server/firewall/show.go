package serverfirewall

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// ShowCommand is the 'server firewall show' command, displaying firewall details
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show server firewall details.",
			"upctl server firewall show 00038afc-d526-4148-af0e-d2f1eeaded9b",
			"upctl server firewall show 00038afc-d526-4148-af0e-d2f1eeaded9b 009d7f4e-99ce-4c78-88f1-e695d4c37743",
			"upctl server firewall show my_server",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
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

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// get rules and server details in parallel
	var wg sync.WaitGroup
	var rules *upcloud.FirewallRules
	var fwRuleErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		rules, fwRuleErr = exec.Firewall().GetFirewallRules(exec.Context(), &request.GetFirewallRulesRequest{ServerUUID: arg})
	}()
	server, err := exec.Server().GetServerDetails(exec.Context(), &request.GetServerDetailsRequest{UUID: arg})
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

	return output.MarshaledWithHumanOutput{
		Value: rules,
		Output: output.Combined{
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
						{Key: "enabled", Title: "Enabled", Value: server.Firewall == "on", Format: format.Boolean},
					}},
				}},
			},
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
