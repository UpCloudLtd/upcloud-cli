package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
)

func ShowCommand(service ServerFirewall) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show server details"),
		service:     service,
	}
}

type showCommand struct {
	*commands.BaseCommand
	service       ServerFirewall
	firewallRules *upcloud.FirewallRules
}

func (s *showCommand) InitCommand() {
	s.ArgCompletion(func(toComplete string) ([]string, cobra.ShellCompDirective) {
		servers, err := s.service.GetServers()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range servers.Servers {
			vals = append(vals, v.UUID, v.Hostname)
		}
		return commands.MatchStringPrefix(vals, toComplete, false),
			cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		// TODO(aakso): implement prompting with readline support
		if len(args) < 1 {
			return nil, fmt.Errorf("server hostname, title or uuid is required")
		}
		var servers []upcloud.Server
		server, err := searchServer(&servers, s.service, args[0], true)
		if err != nil {
			return nil, err
		}
		serverUuid := server.UUID
		var (
			wg        sync.WaitGroup
			fwRuleErr error
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.firewallRules, fwRuleErr = s.service.GetFirewallRules(&request.GetFirewallRulesRequest{ServerUUID: serverUuid})
		}()
		serverDetails, err := s.service.GetServerDetails(&request.GetServerDetailsRequest{UUID: serverUuid})
		if err != nil {
			return nil, err
		}
		wg.Wait()
		if fwRuleErr != nil {
			return nil, fwRuleErr
		}
		return serverDetails, nil
	}
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	srv := out.(*upcloud.ServerDetails)

	dMain := ui.NewDetailsView()
	dMain.SetRowSeparators(true)
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
	sb := &strings.Builder{}

	// Common details
	{
		dCommon := ui.NewDetailsView()
		dCommon.SetRowTransformer(rowTransformer)
		dCommon.AppendRows([]table.Row{
			{"UUID", ui.DefaultUuidColours.Sprint(srv.UUID)},
			{"Title", srv.Title},
			{"Hostname", srv.Hostname},
			{"Plan", srv.Plan},
			{"Zone", srv.Zone},
			{"State", StateColour(srv.State).Sprint(srv.State)},
			{"Tags", strings.Join(srv.Tags, ",")},
			{"License", srv.License},
			{"Metadata", srv.Metadata},
			{"Timezone", srv.Timezone},
			{"Host ID", srv.Host},
		})
		// fmt.Println(dCommon.Render())
		dMain.AppendRow(table.Row{"Common", dCommon.Render()})
	}

	// Storage details
	{
		dStorage := ui.NewDetailsView()
		dStorage.SetRowSpacing(true)
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
		sb.Reset()
		sb.WriteString(tStorage.Render())
		_, _ = fmt.Fprint(sb, "\n\nFlags: B = bootdisk, P = part of plan")
		dStorage.AppendRow(table.Row{"Devices", sb.String()})
		simpleBackup := srv.SimpleBackup
		if simpleBackup == "no" {
			simpleBackup = ui.DefaultBooleanColoursFalse.Sprint(simpleBackup)
		}
		dStorage.AppendRow(table.Row{"Simple Backup", simpleBackup})
		dMain.AppendRow(table.Row{"Storage", dStorage.Render()})
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
				addrs = append(addrs, prefix+ui.DefaultAddressColours.Sprint(addr.Address))
			}
			tNics.AppendRow(table.Row{
				nic.Index,
				nic.Type,
				ui.DefaultUuidColours.Sprint(nic.Network),
				strings.Join(addrs, "\n"),
				strings.Join(flags, " "),
			})
		}
		sb.Reset()
		sb.WriteString(tNics.Render())
		_, _ = fmt.Fprint(sb, "\n\nFlags: S = source IP filtering, B = bootable")
		dNetwork.AppendRow(table.Row{"NICS", sb.String()})
		dNetwork.AppendRow(table.Row{"NIC Model", srv.NICModel})
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
			for _, rule := range s.firewallRules.FirewallRules {
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
			dNetwork.AppendRow(table.Row{"Firewall\nRules", tFw.Render()})
		}
		dMain.AppendRow(table.Row{"Networking", dNetwork.Render()})
	}

	// Remote access
	{
		dRemote := ui.NewDetailsView()
		dRemote.AppendRow(table.Row{"Enabled", formatBool(srv.RemoteAccessEnabled.Bool())})
		if srv.RemoteAccessEnabled.Bool() {
			dRemote.AppendRows([]table.Row{
				{"Type", srv.RemoteAccessType},
				{"Address", fmt.Sprintf("%s:%d", srv.RemoteAccessHost, srv.RemoteAccessPort)},
				{"Password", srv.RemoteAccessPassword},
			})
		}
		dMain.AppendRow(table.Row{"Remote Access", dRemote.Render()})
	}

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, dMain.Render())
	fmt.Fprintln(writer)
	return nil
}
