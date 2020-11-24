package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"
	"io"
)

func ListCommand(service service.Network) commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List networks"),
		service:     service,
	}
}

type listCommand struct {
	*commands.BaseCommand
	service        service.Network
	header         table.Row
	columnKeys     []string
	visibleColumns []string
}

func (s *listCommand) InitCommand() {
	s.header = table.Row{"UUID", "Name", "Router", "Zone", "Type"}
	s.columnKeys = []string{"uuid", "name", "router", "zone", "type"}
	s.visibleColumns = []string{"uuid", "name", "router", "zone", "type"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		ips, err := s.service.GetNetworks()
		if err != nil {
			return nil, err
		}
		return ips, nil
	}
}

func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
	networks := out.(*upcloud.Networks)

	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	for _, n := range networks.Networks {
		t.AppendRow(table.Row{
			n.UUID,
			n.Name,
			n.Router,
			n.Type,
			n.Zone})
	}

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, t.Render())
	fmt.Fprintln(writer)
	return nil
}
