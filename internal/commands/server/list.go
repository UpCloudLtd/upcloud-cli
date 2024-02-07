package server

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/pflag"
)

// ListCommand creates the "server list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List current servers",
			"upctl server list",
			"upctl server list --show-ip-addresses",
			"upctl server list --show-ip-addresses=public",
		),
	}
}

type listServerIpaddresses struct {
	ServerUUID  string
	IPAddresses upcloud.IPAddressSlice
	Error       error
}

type serverWithIPAddress struct {
	upcloud.Server

	IPAddresses upcloud.IPAddressSlice `json:"ip_addresses"`
}

type serversWithIPAddresses struct {
	Servers []serverWithIPAddress `json:"servers"`
}

type listCommand struct {
	*commands.BaseCommand
	showIPAddresses string
}

// InitCommand implements Command.InitCommand
func (ls *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&ls.showIPAddresses, "show-ip-addresses", "none", "Show servers IP addresses of specified access type in the output or all ip addresses if argument value is \"all\" or no argument is specified.")
	flags.Lookup("show-ip-addresses").NoOptDefVal = "all"
	ls.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (ls *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	servers, err := svc.GetServers(exec.Context())
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

		rows = append(rows, output.TableRow{
			s.UUID,
			s.Hostname,
			plan,
			s.Zone,
			s.State,
		})
	}

	columns := []output.TableColumn{
		{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
		{Key: "hostname", Header: "Hostname"},
		{Key: "plan", Header: "Plan"},
		{Key: "zone", Header: "Zone"},
		{Key: "state", Header: "State", Format: format.ServerState},
	}

	if ls.showIPAddresses != "none" {
		serversWithIPs := serversWithIPAddresses{}
		ipaddressMap, err := getIPAddressesByServerUUID(servers, ls.showIPAddresses, exec)
		if err != nil {
			return nil, err
		}

		for i, row := range rows {
			uuid := row[0].(string)

			// Add IPs to the table in human output
			var listIpaddresses upcloud.IPAddressSlice
			if serverIpaddresses, ok := ipaddressMap[uuid]; ok {
				listIpaddresses = append(listIpaddresses, serverIpaddresses.IPAddresses...)
			}
			row = append(row[:3], row[2:]...)
			row[2] = listIpaddresses
			rows[i] = row

			// Add IPs to machine output
			serversWithIPs.Servers = append(serversWithIPs.Servers, serverWithIPAddress{
				Server:      servers.Servers[i],
				IPAddresses: listIpaddresses,
			})
		}
		columns = append(columns[:3], columns[2:]...)
		columns[2] = output.TableColumn{
			Key:    "ip_addresses",
			Header: "IP addresses",
			Format: formatListIPAddresses,
		}

		return output.MarshaledWithHumanOutput{
			Value: serversWithIPs,
			Output: output.Table{
				Columns: columns,
				Rows:    rows,
			},
		}, nil
	}

	return output.MarshaledWithHumanOutput{
		Value: servers,
		Output: output.Table{
			Columns: columns,
			Rows:    rows,
		},
	}, nil
}

// getIPAddressesByServerUUID returns IP addresses grouped by server UUID. This function will be removed when server end-point response includes IP addresses.
func getIPAddressesByServerUUID(servers *upcloud.Servers, accessType string, exec commands.Executor) (map[string]listServerIpaddresses, error) {
	returnChan := make(chan listServerIpaddresses)
	var wg sync.WaitGroup

	for _, server := range servers.Servers {
		wg.Add(1)
		go func(server upcloud.Server) {
			defer wg.Done()
			ipaddresses, err := getServerIPAddresses(server.UUID, accessType, exec)
			returnChan <- listServerIpaddresses{
				ServerUUID:  server.UUID,
				IPAddresses: ipaddresses,
				Error:       err,
			}
		}(server)
	}

	go func() {
		wg.Wait()
		close(returnChan)
	}()

	ipaddressMap := make(map[string]listServerIpaddresses)
	for response := range returnChan {
		ipaddressMap[response.ServerUUID] = response
	}

	return ipaddressMap, nil
}

func getServerIPAddresses(uuid, accessType string, exec commands.Executor) (upcloud.IPAddressSlice, error) {
	server, err := exec.All().GetServerNetworks(exec.Context(), &request.GetServerNetworksRequest{ServerUUID: uuid})
	if err != nil {
		return nil, err
	}

	var ipaddresses upcloud.IPAddressSlice
	for _, iface := range server.Interfaces {
		for _, ipa := range iface.IPAddresses {
			if accessType == "all" || iface.Type == accessType {
				ipa.Access = iface.Type
				ipaddresses = append(ipaddresses, ipa)
			}
		}
	}

	sort.Slice(ipaddresses, func(i, j int) bool {
		accessMap := map[string]int{
			"public":  3,
			"private": 2,
			"utility": 1,
		}
		floatingMap := map[bool]int{
			true:  1,
			false: 0,
		}

		if accessMap[ipaddresses[i].Access] != accessMap[ipaddresses[j].Access] {
			return accessMap[ipaddresses[i].Access] > accessMap[ipaddresses[j].Access]
		}

		return floatingMap[ipaddresses[i].Floating.Bool()] > floatingMap[ipaddresses[j].Floating.Bool()]
	})

	return ipaddresses, nil
}

func formatListIPAddresses(val interface{}) (text.Colors, string, error) {
	ipaddresses, ok := val.(upcloud.IPAddressSlice)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse IP addresses from %T, expected upcloud.IPAddressSlice", val)
	}

	var rows []string
	for _, ipa := range ipaddresses {
		var floating string
		if ipa.Floating.Bool() {
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
