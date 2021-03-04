package serverfirewall

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/UpCloudLtd/cli/internal/commands/server"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ShowCommand(serverSvc service.Server, firewallSvc service.Firewall) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show server details."),
		serverSvc:   serverSvc,
		firewallSvc: firewallSvc,
	}
}

type showCommand struct {
	*commands.BaseCommand
	serverSvc   service.Server
	firewallSvc service.Firewall
}

type commandResponseHolder struct {
	serverDetails *upcloud.ServerDetails
	firewallRules *upcloud.FirewallRules
}

// MarshalJSON implements json.Marshaler
func (c *commandResponseHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.serverDetails)
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(server.GetServerArgumentCompletionFunction(s.serverSvc))
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		// TODO(aakso): implement prompting with readline support
		if len(args) != 1 {
			return nil, fmt.Errorf("one server hostname, title or uuid is required")
		}
		serverUUIDs, err := server.SearchAllServers(args, s.serverSvc, true)
		if err != nil {
			return nil, err
		}
		if len(serverUUIDs) != 1 {
			return nil, fmt.Errorf("server not found")
		}
		serverUUID := serverUUIDs[0]
		var (
			wg        sync.WaitGroup
			fwRuleErr error
		)
		wg.Add(1)
		var firewallRules *upcloud.FirewallRules
		go func() {
			defer wg.Done()
			firewallRules, fwRuleErr = s.firewallSvc.GetFirewallRules(&request.GetFirewallRulesRequest{ServerUUID: serverUUID})
		}()
		serverDetails, err := s.serverSvc.GetServerDetails(&request.GetServerDetailsRequest{UUID: serverUUID})
		if err != nil {
			return nil, err
		}
		wg.Wait()
		if fwRuleErr != nil {
			return nil, fwRuleErr
		}
		return &commandResponseHolder{serverDetails, firewallRules}, nil
	}
}

// HandleOutput implements Command.HandleOutput
func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	resp := out.(*commandResponseHolder)
	srv := resp.serverDetails
	firewallRules := resp.firewallRules

	formatBool := func(v bool) interface{} {
		if v {
			return ui.DefaultBooleanColoursTrue.Sprint("yes")
		}
		return ui.DefaultBooleanColoursFalse.Sprint("no")
	}

	l := ui.NewListLayout(ui.ListLayoutDefault)

	{
		dNetwork := ui.NewDetailsView()
		dNetwork.SetRowSpacing(true)
		formatMatch := func(start, stop, portStart, portStop string) string {
			var sb strings.Builder
			ipStart := net.ParseIP(start)
			ipStop := net.ParseIP(stop)
			if ipStart != nil {
				if ipStart.Equal(ipStop) {
					sb.WriteString(ui.DefaultAddressColours.Sprint(ipStart))
				} else {
					sb.WriteString(ui.DefaultAddressColours.Sprintf("%s →\n%s", ipStart, ipStop))
				}
			}
			if portStart != "" {
				if ipStart != nil {
					sb.WriteString("\n")
				}
				if portStart == portStop {
					sb.WriteString(fmt.Sprintf("port: %s", portStart))
				} else {
					sb.WriteString(fmt.Sprintf("port: %s → %s", portStart, portStop))
				}
			}
			return sb.String()
		}
		formatProto := func(family, proto, icmptype string) string {
			var sb strings.Builder
			sb.WriteString(family)
			if proto == "" {
				return sb.String()
			}
			sb.WriteString(fmt.Sprintf("/%s", proto))
			if icmptype == "" {
				return sb.String()
			}
			sb.WriteString(fmt.Sprintf("/%s", icmptype))
			return sb.String()
		}
		tFw := ui.NewDataTable(
			"#",
			"Action",
			"Source",
			"Destination",
			"Dir",
			"Proto",
		)
		tFw.SetColumnConfig("Source", table.ColumnConfig{WidthMax: 27})
		tFw.SetColumnConfig("Destination", table.ColumnConfig{WidthMax: 27})

		for _, rule := range firewallRules.FirewallRules {
			actColour := text.FgHiGreen
			if rule.Action == "drop" {
				actColour = text.FgHiRed
			}
			tFw.Append(table.Row{
				rule.Position,
				actColour.Sprint(rule.Action),
				formatMatch(
					rule.SourceAddressStart,
					rule.SourceAddressEnd,
					rule.SourcePortStart,
					rule.SourcePortEnd),
				formatMatch(
					rule.DestinationAddressStart,
					rule.DestinationAddressEnd,
					rule.DestinationPortStart,
					rule.DestinationPortEnd),
				rule.Direction,
				formatProto(
					rule.Family,
					rule.Protocol,
					rule.ICMPType),
			})
		}

		l.AppendSection("Firewall Rules:", ui.WrapWithListLayout(tFw.Render(), ui.ListLayoutNestedTable).Render())
	}

	// Firewall status
	fwEnabled := srv.Firewall == "on"
	{
		dRemote := ui.NewDetailsView()
		dRemote.Append(table.Row{"Enabled:", formatBool(fwEnabled)})
		l.AppendSection(dRemote.Render())
	}

	_, _ = fmt.Fprintln(writer, l.Render())
	return nil

}
