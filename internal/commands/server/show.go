package server

import (
	"fmt"
	"strings"
	"sync"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// ShowCommand creates the "server show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show server details",
			"upctl server show 21aeb3b7-cd89-4123-a376-559b0e75be8b",
			"upctl server show 21aeb3b7-cd89-4123-a376-559b0e75be8b 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0",
			"upctl server show myhostname",
			"upctl server show my_server1 my_server2",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
}

func (s *showCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	var (
		wg        sync.WaitGroup
		fwRuleErr error
	)

	serverSvc := exec.Server()
	firewallSvc := exec.Firewall()

	wg.Add(1)
	var firewallRules *upcloud.FirewallRules
	go func() {
		defer wg.Done()
		firewallRules, fwRuleErr = firewallSvc.GetFirewallRules(&request.GetFirewallRulesRequest{ServerUUID: uuid})
	}()

	server, err := serverSvc.GetServerDetails(&request.GetServerDetailsRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	wg.Wait()
	if fwRuleErr != nil {
		return nil, fwRuleErr
	}

	planOutput := server.Plan
	if planOutput == "custom" {
		memory := server.MemoryAmount / 1024
		planOutput = fmt.Sprintf("%dxCPU-%dGB (custom)", server.CoreNumber, memory)
	}

	// Storage details
	storageRows := []output.TableRow{}
	for _, storage := range server.StorageDevices {
		var flags []string
		if storage.PartOfPlan == "yes" {
			flags = append(flags, "P")
		}
		if storage.BootDisk == 1 {
			flags = append(flags, "B")
		}

		storageRows = append(storageRows, output.TableRow{
			storage.UUID,
			storage.Title,
			storage.Type,
			storage.Address,
			storage.Size,
			strings.Join(flags, " "),
		})
	}
	// Network details
	nicRows := []output.TableRow{}
	for _, nic := range server.Networking.Interfaces {
		var flags []string
		if nic.SourceIPFiltering.Bool() {
			flags = append(flags, "S")
		}
		if nic.Bootable.Bool() {
			flags = append(flags, "B")
		}

		var addrs []string
		for _, addr := range nic.IPAddresses {
			var floating string
			if addr.Floating.Bool() {
				floating = " (f)"
			}

			addrs = append(
				addrs,
				fmt.Sprintf(
					"%v: %s%s",
					addr.Family,
					ui.DefaultAddressColours.Sprint(addr.Address),
					floating),
			)
		}

		nicRows = append(nicRows, output.TableRow{
			nic.Index,
			nic.Type,
			strings.Join(addrs, "\n"),
			nic.MAC,
			nic.Network,
			strings.Join(flags, " "),
		})
	}

	combined := output.Combined{
		output.CombinedSection{
			Contents: output.Details{
				Sections: []output.DetailSection{
					{
						Title: "Common",
						Rows: []output.DetailRow{
							{Title: "UUID:", Key: "uuid", Value: server.UUID, Colour: ui.DefaultUUUIDColours},
							{Title: "Hostname:", Key: "hostname", Value: server.Hostname},
							{Title: "Title:", Key: "title", Value: server.Title},
							{Title: "Plan:", Key: "plan", Value: planOutput},
							{Title: "Zone:", Key: "zone", Value: server.Zone},
							{Title: "State:", Key: "state", Value: server.State, Colour: commands.ServerStateColour(server.State)},
							{Title: "Simple Backup:", Key: "simple_backup", Value: server.SimpleBackup},
							{Title: "Licence:", Key: "licence", Value: server.License},
							{Title: "Metadata:", Key: "metadata", Value: server.Metadata.String()},
							{Title: "Timezone:", Key: "timezone", Value: server.Timezone},
							{Title: "Host ID:", Key: "host_id", Value: server.Host},
							{Title: "Tags:", Key: "tags", Value: strings.Join(server.Tags, ",")},
						},
					},
				},
			},
		},
		output.CombinedSection{
			Key:   "storage",
			Title: "Storage: (Flags: B = bootdisk, P = part of plan)",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
					{Key: "title", Header: "Title"},
					{Key: "type", Header: "Type"},
					{Key: "address", Header: "Address"},
					{Key: "size", Header: "Size (GiB)"},
					{Key: "flags", Header: "Flags"},
				},
				Rows: storageRows,
			},
		},
		output.CombinedSection{
			Key:   "nics",
			Title: "NICs: (Flags: S = source IP filtering, B = bootable)",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "id", Header: "#"},
					{Key: "type", Header: "Type"},
					{Key: "ip_address", Header: "IP Address"},
					{Key: "mac_address", Header: "MAC Address"},
					{Key: "network", Header: "Network", Colour: ui.DefaultUUUIDColours},
					{Key: "flags", Header: "Flags"},
				},
				Rows: nicRows,
			},
		},
	}

	// Firewall rules
	if server.Firewall == "on" {
		fwRows := []output.TableRow{}
		for _, rule := range firewallRules.FirewallRules {
			fwRows = append(fwRows, output.TableRow{
				rule.Position,
				rule.Direction,
				rule.Action,
				ui.FormatRange(
					rule.SourceAddressStart,
					rule.SourceAddressEnd,
				),
				ui.FormatRange(
					rule.DestinationAddressStart,
					rule.DestinationAddressEnd,
				),
				ui.FormatRange(
					rule.SourcePortStart,
					rule.SourcePortEnd,
				),
				ui.FormatRange(
					rule.DestinationPortStart,
					rule.DestinationPortEnd,
				),
				ui.ConcatStrings(rule.Family, rule.Protocol, rule.ICMPType),
			})
		}
		combined = append(combined, output.CombinedSection{
			Key:   "firewall",
			Title: "Firewall Rules:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "position", Header: "#"},
					{Key: "direction", Header: "Direction"},
					{Key: "action", Header: "Action"},
					{Key: "src_ipaddress", Header: "Src IPAddress", Colour: ui.DefaultAddressColours},
					{Key: "dest_ipaddress", Header: "Dest IPAddress", Colour: ui.DefaultAddressColours},
					{Key: "src_port", Header: "Src Port"},
					{Key: "dest_port", Header: "Dest Port"},
					{Key: "protocol", Header: "Protocol"},
				},
				Rows: fwRows,
			},
		})
	}

	// Remote access
	if server.RemoteAccessEnabled.Bool() {
		combined = append(combined, output.CombinedSection{
			Contents: output.Details{
				Sections: []output.DetailSection{
					{
						Title: "Remote Access",
						Rows: []output.DetailRow{
							{Title: "Type:", Key: "type", Value: server.RemoteAccessType},
							{Title: "Host:", Key: "host", Value: server.RemoteAccessHost},
							{Title: "Port:", Key: "port", Value: server.RemoteAccessPort},
							{Title: "Password:", Key: "password", Value: server.RemoteAccessPassword},
						},
					},
				},
			},
		})
	}

	return combined, nil
}
