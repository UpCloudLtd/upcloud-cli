package network

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
	zone           string
}

func (s *listCommand) InitCommand() {
	s.header = table.Row{"UUID", "Name", "Router", "Type", "Zone"}
	s.columnKeys = []string{"uuid", "name", "router", "type", "zone"}
	s.visibleColumns = []string{"uuid", "name", "router", "type", "zone"}
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.zone, "zone", "", "Filters for given zone")
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if s.zone != "" {
			filtered, err := s.service.GetNetworksInZone(&request.GetNetworksInZoneRequest{Zone: s.zone})
			if err != nil {
				return nil, err
			}
			return filtered, nil
		}

		networks, err := s.service.GetNetworks()
		if err != nil {
			return nil, err
		}

		return networks, nil
	}
}

func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
	networks := out.(*upcloud.Networks)

	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	for _, n := range networks.Networks {
		t.AppendRow(table.Row{
			ui.DefaultUuidColours.Sprint(n.UUID),
			n.Name,
			n.Router,
			n.Type,
			n.Zone,
		})
	}

	return t.Paginate(writer)
}
