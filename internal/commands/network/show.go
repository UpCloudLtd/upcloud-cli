package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"strings"
)

func ShowCommand(networkSvc service.Network, serverSvc service.Server) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show current network"),
		networkSvc:  networkSvc,
		serverSvc:   serverSvc,
	}
}

type showCommand struct {
	*commands.BaseCommand
	networkSvc service.Network
	serverSvc  service.Server
}

type networkWithServers struct {
	network *upcloud.Network
	servers []upcloud.Server
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("one network uuid or name is required")
		}
		n, err := SearchNetwork(args[0], s.networkSvc)
		if err != nil {
			return nil, err
		}

		var servers []upcloud.Server
		for _, networkServer := range n.Servers {
			svr, err := server.SearchServer(&servers, s.serverSvc, networkServer.ServerUUID, true)
			if err != nil {
				return nil, err
			}
			servers = append(servers, *svr[0])
		}

		return networkWithServers{network: n, servers: servers}, nil
	}
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	networkWithServers := out.(networkWithServers)
	n := networkWithServers.network
	servers := networkWithServers.servers

	l := ui.NewListLayout(ui.ListLayoutDefault)

	dCommon := ui.NewDetailsView()
	dCommon.AppendRows([]table.Row{
		{"UUID:", ui.DefaultUuidColours.Sprint(n.UUID)},
		{"Name:", n.Name},
		{"Router:", n.Router},
		{"Type:", n.Type},
		{"Zone:", n.Zone},
	})
	l.AppendSection("Common", dCommon.Render())

	if len(n.IPNetworks) > 0 {
		tIPNetwork := ui.NewDataTable("Address", "Family", "DHCP", "DHCP Def Router", "DHCP DNS")
		for _, nip := range n.IPNetworks {
			tIPNetwork.AppendRow(table.Row{
				ui.DefaultAddressColours.Sprint(nip.Address),
				nip.Family,
				ui.FormatBool(nip.DHCP.Bool()),
				ui.FormatBool(nip.DHCPDefaultRoute.Bool()),
				strings.Join(nip.DHCPDns, " "),
			})
		}
		l.AppendSection("IP Networks:", tIPNetwork.Render())
	} else {
		l.AppendSection("IP Networks:", "(no ip network found)")
	}

	if len(servers) > 0 {
		tServers := ui.NewDataTable("Title (UUID)", "Hostname", "State")

		for _, s := range servers {
			tServers.AppendRow(table.Row{
				fmt.Sprintf("%s\n(%s)", s.Title, ui.DefaultUuidColours.Sprint(s.UUID)),
				s.Title,
				s.Hostname,
				server.StateColour(s.State).Sprint(s.State),
			})
		}
		l.AppendSection("Servers:", ui.WrapWithListLayout(tServers.Render(), ui.ListLayoutNestedTable).Render())
	} else {
		l.AppendSection("Servers:", "(no servers using this storage)")
	}

	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
