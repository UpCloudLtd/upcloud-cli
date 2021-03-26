package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowCommand creates the "server show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show server details"),
	}
}

type showCommand struct {
	*commands.BaseCommand
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
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		// TODO(aakso): implement prompting with readline support
		if len(args) != 1 {
			return nil, fmt.Errorf("one server hostname, title or uuid is required")
		}
		serverSvc := s.Config().Service.(service.Server)
		firewallSvc := s.Config().Service.(service.Firewall)
		serverUUIDs, err := SearchAllServers(args, serverSvc, true)
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
			firewallRules, fwRuleErr = firewallSvc.GetFirewallRules(&request.GetFirewallRulesRequest{ServerUUID: serverUUID})
		}()
		serverDetails, err := serverSvc.GetServerDetails(&request.GetServerDetailsRequest{UUID: serverUUID})
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
	rowTransformer := func(row table.Row) table.Row {
		if v, ok := row[len(row)-1].(upcloud.Boolean); ok {
			row[len(row)-1] = formatBool(v.Bool())
		}
		return row
	}

	l := ui.NewListLayout(ui.ListLayoutDefault)
	{
		dCommon := ui.NewDetailsView()
		dCommon.SetRowTransformer(rowTransformer)
		planOutput := srv.Plan
		if planOutput == "custom" {
			memory := srv.MemoryAmount / 1024
			planOutput = fmt.Sprintf("Custom (%dxCPU, %dGB)", srv.CoreNumber, memory)
		}
		dCommon.Append(
			table.Row{"UUID:", ui.DefaultUUUIDColours.Sprint(srv.UUID)},
			table.Row{"Title:", srv.Title},
			table.Row{"Hostname:", srv.Hostname},
			table.Row{"Plan:", planOutput},
			table.Row{"Zone:", srv.Zone},
			table.Row{"State:", commands.StateColour(srv.State).Sprint(srv.State)},
			table.Row{"Tags:", strings.Join(srv.Tags, ",")},
			table.Row{"Licence:", srv.License},
			table.Row{"Metadata:", srv.Metadata},
			table.Row{"Timezone:", srv.Timezone},
			table.Row{"Host ID:", srv.Host},
		)
		l.AppendSection("Common:", dCommon.Render())
	}

	// Storage details
	{
		tStorage := ui.NewDataTable("Title (UUID)", "Type", "Address", "Size (GiB)", "Flags")
		for _, storage := range srv.StorageDevices {
			var flags []string
			if storage.PartOfPlan == "yes" {
				flags = append(flags, "P")
			}
			if storage.BootDisk == 1 {
				flags = append(flags, "B")
			}
			tStorage.Append(table.Row{
				fmt.Sprintf("%s\n(%s)", storage.Title, ui.DefaultUUUIDColours.Sprint(storage.UUID)),
				storage.Type,
				storage.Address,
				storage.Size,
				strings.Join(flags, " "),
			})
		}

		simpleBackup := srv.SimpleBackup
		if simpleBackup == "no" {
			simpleBackup = ui.DefaultBooleanColoursFalse.Sprint(simpleBackup)
		}

		dStorage := ui.NewListLayout(ui.ListLayoutNestedTable)

		dBackup := ui.NewDetailsView()
		dBackup.SetRowTransformer(rowTransformer)
		dBackup.Append(table.Row{"Simple Backup:", simpleBackup})
		dStorage.AppendSectionWithNote("Devices:", tStorage.Render(), "(Flags: B = bootdisk, P = part of plan)")

		l.AppendSection("Storage:", dBackup.Render(), dStorage.Render())
	}

	// Network details
	{
		dNetwork := ui.NewDetailsView()
		dNetwork.SetRowSpacing(true)
		tNics := ui.NewDataTable("#", "Type", "Network", "Addresses", "Flags")
		tNics.SetColumnConfig("#", table.ColumnConfig{WidthMax: 2})
		tNics.SetColumnConfig("Network", table.ColumnConfig{WidthMax: 19})
		for _, nic := range srv.Networking.Interfaces {
			var flags []string
			if nic.SourceIPFiltering.Bool() {
				flags = append(flags, "S")
			}
			if nic.Bootable.Bool() {
				flags = append(flags, "B")
			}
			var addrs []string
			addrs = append(addrs, "MAC:  "+nic.MAC)
			for _, addr := range nic.IPAddresses {
				prefix := "IPv4: "
				if addr.Family == "IPv6" {
					prefix = "IPv6: "
				}
				var floating string
				if addr.Floating.Bool() {
					floating = " (f)"
				}
				addrs = append(addrs, prefix+ui.DefaultAddressColours.Sprint(addr.Address)+floating)
			}
			tNics.Append(table.Row{
				nic.Index,
				nic.Type,
				ui.DefaultUUUIDColours.Sprint(nic.Network),
				strings.Join(addrs, "\n"),
				strings.Join(flags, " "),
			})
		}

		dNICs := ui.NewListLayout(ui.ListLayoutNestedTable)
		dNICs.AppendSectionWithNote("NICS:", tNics.Render(), "(Flags: S = source IP filtering, B = bootable)")

		l.AppendSection("Networking:", dNICs.Render())

		fwEnabled := srv.Firewall == "on"
		dNetwork.Append(table.Row{"Firewall", formatBool(fwEnabled)})
		if fwEnabled {
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
	}

	// Remote access
	{
		dRemote := ui.NewDetailsView()
		dRemote.Append(table.Row{"Enabled", formatBool(srv.RemoteAccessEnabled.Bool())})
		if srv.RemoteAccessEnabled.Bool() {
			dRemote.Append(
				table.Row{"Type:", srv.RemoteAccessType},
				table.Row{"Address:", fmt.Sprintf("%s:%d", srv.RemoteAccessHost, srv.RemoteAccessPort)},
				table.Row{"Password:", srv.RemoteAccessPassword},
			)
		}
		l.AppendSection("Remote Access:", dRemote.Render())
	}

	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
