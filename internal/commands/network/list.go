package network

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"
)

// ListCommand creates the "network list" command
func ListCommand(service service.Network) commands.NewCommand {
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
func (s *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.zone, "zone", "", "Filters for given zone")
	flags.StringVar(&s.networkType, "type", "", "Filters for given type")
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

// Execute implements command.NewCommand
func (s *listCommand) Execute(exec commands.Executor, args []string) (output.Command, error) {
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

	if s.networkType != "" {
		var filtered []upcloud.Network
		for _, n := range networks.Networks {
			if n.Type == s.networkType {
				filtered = append(filtered, n)
			}
		}
		networks = &upcloud.Networks{Networks: filtered}
	}
	var rows []output.TableRow
	for _, n := range networks.Networks {
		rows = append(rows, output.TableRow{
			// TODO: reimplement
			// ui.DefaultUUUIDColours.Sprint(n.UUID),
			n.UUID,
			n.Name,
			n.Router,
			n.Type,
			n.Zone,
		})
	}
	return output.Table{
		Columns: []output.TableColumn{
			{Key: "uuid", Header: "UUID"},
			{Key: "name", Header: "Name"},
			{Key: "router", Header: "Router"},
			{Key: "type", Header: "Type"},
			{Key: "zone", Header: "Zone"},
		},
		Rows: rows,
	}, nil
}

// NewParent implements command.NewCommand
func (s *listCommand) NewParent() commands.NewCommand {
	return s.Parent().(commands.NewCommand)
}
