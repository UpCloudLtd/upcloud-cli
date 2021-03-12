package network

import (
	"encoding/json"
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

// ShowCommand creates the "network show" command
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

func (s *showCommand) InitCommand() error {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getArgCompFn(s.networkSvc))

	return nil
}

type networkWithServers struct {
	network *upcloud.Network
	servers []*upcloud.Server
}

func (c *networkWithServers) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.network)
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("one network uuid or name is required")
		}
		n, err := SearchUniqueNetwork(args[0], s.networkSvc)
		if err != nil {
			return nil, err
		}

		var servers []*upcloud.Server
		for _, networkServer := range n.Servers {
			svr, err := server.SearchSingleServer(networkServer.ServerUUID, s.serverSvc)
			if err != nil {
				return nil, err
			}
			servers = append(servers, svr)
		}

		return &networkWithServers{network: n, servers: servers}, nil
	}
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	networkWithServers := out.(*networkWithServers)
	n := networkWithServers.network
	servers := networkWithServers.servers

	l := ui.NewListLayout(ui.ListLayoutDefault)

	dCommon := ui.NewDetailsView()
	dCommon.Append(
		table.Row{"UUID:", ui.DefaultUUUIDColours.Sprint(n.UUID)},
		table.Row{"Name:", n.Name},
		table.Row{"Router:", n.Router},
		table.Row{"Type:", n.Type},
		table.Row{"Zone:", n.Zone},
	)
	l.AppendSection("Common", dCommon.Render())

	if len(n.IPNetworks) > 0 {
		tIPNetwork := ui.NewDataTable("Address", "Family", "DHCP", "DHCP Def Router", "DHCP DNS")
		for _, nip := range n.IPNetworks {
			tIPNetwork.Append(table.Row{
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
		tServers := ui.NewDataTable("UUID", "Title", "Hostname", "State")

		for _, s := range servers {
			tServers.Append(table.Row{
				ui.DefaultUUUIDColours.Sprint(s.UUID),
				s.Title,
				s.Hostname,
				commands.StateColour(s.State).Sprint(s.State),
			})
		}
		l.AppendSection("Servers:", ui.WrapWithListLayout(tServers.Render(), ui.ListLayoutNestedTable).Render())
	} else {
		l.AppendSection("Servers:", "(no servers assigned to this network)")
	}

	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
