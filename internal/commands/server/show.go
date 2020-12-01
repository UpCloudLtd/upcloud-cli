package server

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ShowCommand(serverSvc service.Server, firewallSvc service.Firewall) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show server details"),
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

func (c *commandResponseHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.serverDetails)
}

func (s *showCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.serverSvc))
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		// TODO(aakso): implement prompting with readline support
		if len(args) != 1 {
			return nil, fmt.Errorf("one server hostname, title or uuid is required")
		}
		serverUuids, err := SearchAllServers(args, s.serverSvc, true)
		if err != nil {
			return nil, err
		}
		if len(serverUuids) != 1 {
			return nil, fmt.Errorf("server not found")
		}
		serverUuid := serverUuids[0]
		var (
			wg        sync.WaitGroup
			fwRuleErr error
		)
		wg.Add(1)
		var firewallRules *upcloud.FirewallRules
		go func() {
			defer wg.Done()
			firewallRules, fwRuleErr = s.firewallSvc.GetFirewallRules(&request.GetFirewallRulesRequest{ServerUUID: serverUuid})
		}()
		serverDetails, err := s.serverSvc.GetServerDetails(&request.GetServerDetailsRequest{UUID: serverUuid})
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
		dCommon.AppendRows([]table.Row{
			{"UUID:", ui.DefaultUuidColours.Sprint(srv.UUID)},
			{"Title:", srv.Title},
			{"Hostname:", srv.Hostname},
			{"Plan:", srv.Plan},
			{"Zone:", srv.Zone},
			{"State:", commands.StateColour(srv.State).Sprint(srv.State)},
			{"Tags:", strings.Join(srv.Tags, ",")},
			{"Licence:", srv.License},
			{"Metadata:", srv.Metadata},
			{"Timezone:", srv.Timezone},
			{"Host ID:", srv.Host},
		})
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
			tStorage.AppendRow(table.Row{
				fmt.Sprintf("%s\n(%s)", storage.Title, ui.DefaultUuidColours.Sprint(storage.UUID)),
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
		dBackup.AppendRows([]table.Row{
			{"Simple Backup:", simpleBackup},
		})
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
					floating = "(f) "
				}
				addrs = append(addrs, floating+prefix+ui.DefaultAddressColours.Sprint(addr.Address))
			}
			tNics.AppendRow(table.Row{
				nic.Index,
				nic.Type,
				ui.DefaultUuidColours.Sprint(nic.Network),
				strings.Join(addrs, "\n"),
				strings.Join(flags, " "),
			})
		}

		dNICs := ui.NewListLayout(ui.ListLayoutNestedTable)
		dNICs.AppendSectionWithNote("NICS:", tNics.Render(), "(Flags: S = source IP filtering, B = bootable)")

		l.AppendSection("Networking:", dNICs.Render())

		fwEnabled := srv.Firewall == "on"
		dNetwork.AppendRow(table.Row{"Firewall", formatBool(fwEnabled)})
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
				tFw.AppendRow(table.Row{
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
		dRemote.AppendRow(table.Row{"Enabled", formatBool(srv.RemoteAccessEnabled.Bool())})
		if srv.RemoteAccessEnabled.Bool() {
			dRemote.AppendRows([]table.Row{
				{"Type:", srv.RemoteAccessType},
				{"Address:", fmt.Sprintf("%s:%d", srv.RemoteAccessHost, srv.RemoteAccessPort)},
				{"Password:", srv.RemoteAccessPassword},
			})
		}
		l.AppendSection("Remote Access:", dRemote.Render())
	}

	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
