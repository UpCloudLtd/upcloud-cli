package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"strings"
)

func ShowCommand(service service.Network) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show current network"),
		service:     service,
	}
}

type showCommand struct {
	*commands.BaseCommand
	service service.Network
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("one network uuid or name is required")
		}
		ip, err := searchNetwork(args[0], s.service)
		if err != nil {
			return nil, err
		}
		return ip, nil
	}
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	n := out.(*upcloud.Network)

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

	tIPNetwork := ui.NewDataTable("Address", "Family", "DHCP", "DHCP Def Router", "DHCP DNS")
	for _, nip := range n.IPNetworks {
		tIPNetwork.AppendRow(table.Row{
			ui.DefaultAddressColours.Sprint(nip.Address),
			nip.Family,
			nip.DHCP == 1,
			nip.DHCPDefaultRoute == 1,
			strings.Join(nip.DHCPDns, " "),
		})
	}
	l.AppendSection("IP Networks:", tIPNetwork.Render())

	tServer := ui.NewDataTable("UUID", "Title")
	for _, server := range n.Servers {
		tIPNetwork.AppendRow(table.Row{
			ui.DefaultUuidColours.Sprint(server.ServerUUID),
			server.ServerTitle,
		})
	}
	l.AppendSection("Server:", tServer.Render())

	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
