package server

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/service"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/pflag"
)

// ListCommand creates the "server list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current servers", "upctl server list"),
	}
}

type listIPAddress struct {
	Access   string `json:"access"`
	Address  string `json:"address"`
	Floating bool   `json:"floating"`
}

type listCommand struct {
	*commands.BaseCommand
	showIPAddresses config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (ls *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &ls.showIPAddresses, "show-ip-addresses", false, "Show IP addresses of the servers in the output.")
	ls.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (ls *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	servers, err := svc.GetServers()
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, s := range servers.Servers {
		plan := s.Plan
		if plan == customPlan {
			memory := s.MemoryAmount / 1024
			plan = fmt.Sprintf("%dxCPU-%dGB (custom)", s.CoreNumber, memory)
		}

		coloredState := commands.ServerStateColour(s.State).Sprint(s.State)

		rows = append(rows, output.TableRow{
			s.UUID,
			s.Hostname,
			plan,
			s.Zone,
			coloredState,
		})
	}

	columns := []output.TableColumn{
		{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
		{Key: "hostname", Header: "Hostname"},
		{Key: "plan", Header: "Plan"},
		{Key: "zone", Header: "Zone"},
		{Key: "state", Header: "State"},
	}

	if ls.showIPAddresses.Value() {
		ipaddressMap, err := getIPAddressesByServerUUID(svc)
		if err != nil {
			return nil, err
		}

		for i, row := range rows {
			uuid := row[0].(string)

			var listIpaddresses []listIPAddress
			if apiIpaddresses, ok := ipaddressMap[uuid]; ok {
				for _, ipa := range apiIpaddresses {
					listIpaddresses = append(listIpaddresses, listIPAddress{
						Access:   ipa.Access,
						Address:  ipa.Address,
						Floating: ipa.Floating.Bool(),
					})
				}
			}
			row = append(row[:3], row[2:]...)
			row[2] = listIpaddresses
			rows[i] = row
		}
		columns = append(columns[:3], columns[2:]...)
		columns[2] = output.TableColumn{
			Key:    "ip_addresses",
			Header: "IP addresses",
			Format: formatListIPAddresses,
		}
	}

	return output.Table{
		Columns: columns,
		Rows:    rows,
	}, nil
}

// getIPAddressesByServerUUID returns IP addresses grouped by server UUID. This function will be removed when server end-point response includes IP addresses.
func getIPAddressesByServerUUID(svc service.AllServices) (map[string][]upcloud.IPAddress, error) {
	ipaddresses, err := svc.GetIPAddresses()
	if err != nil {
		return nil, err
	}

	ipaddressMap := make(map[string][]upcloud.IPAddress)
	for _, ipaddress := range ipaddresses.IPAddresses {
		current, ok := ipaddressMap[ipaddress.ServerUUID]
		if ok {
			ipaddressMap[ipaddress.ServerUUID] = append(current, ipaddress)
		} else {
			ipaddressMap[ipaddress.ServerUUID] = []upcloud.IPAddress{ipaddress}
		}
	}

	return ipaddressMap, nil
}

func formatListIPAddresses(val interface{}) (text.Colors, string, error) {
	ipaddresses, ok := val.([]listIPAddress)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse IP addresses from %T, expected []listIPAddress", val)
	}

	var rows []string
	for _, ipa := range ipaddresses {
		var floating string
		if ipa.Floating {
			floating = " (f)"
		}

		rows = append(rows, fmt.Sprintf(
			"%s: %s%s",
			ipa.Access,
			ui.DefaultAddressColours.Sprint(ipa.Address),
			floating,
		))
	}

	return nil, strings.Join(rows, ",\n"), nil
}
