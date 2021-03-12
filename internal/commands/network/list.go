package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"
	"io"
)

// ListCommand creates the "network list" command
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
	networkType    string
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() error {
	s.header = table.Row{"UUID", "Name", "Router", "Type", "Zone"}
	s.columnKeys = []string{"uuid", "name", "router", "type", "zone"}
	s.visibleColumns = []string{"uuid", "name", "router", "type", "zone"}
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.zone, "zone", "", "Filters for given zone")
	flags.StringVar(&s.networkType, "type", "", "Filters for given type")
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)

	return nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		var filtered []upcloud.Network

		var networks *upcloud.Networks
		var err error
		if s.zone != "" {
			networks, err = s.service.GetNetworksInZone(&request.GetNetworksInZoneRequest{Zone: s.zone})
		} else {
			networks, err = s.service.GetNetworks()
		}
		if err != nil {
			return nil, err
		}

		if s.networkType == "" {
			return networks, nil
		}

		for _, n := range networks.Networks {
			if n.Type == s.networkType {
				filtered = append(filtered, n)
			}
		}
		return &upcloud.Networks{Networks: filtered}, nil
	}
}

// HandleOutput implements Command.HandleOutput
func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
	networks := out.(*upcloud.Networks)

	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	for _, n := range networks.Networks {
		t.Append(table.Row{
			ui.DefaultUUUIDColours.Sprint(n.UUID),
			n.Name,
			n.Router,
			n.Type,
			n.Zone,
		})
	}

	_, _ = fmt.Fprintln(writer, t.Render())
	return nil
}
