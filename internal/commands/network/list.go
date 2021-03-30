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
		BaseCommand: commands.New("list", "List networks, by default private networks only"),
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
	all            bool
	public         bool
	utility        bool
	private        bool
}

func (s *listCommand) MakeExecutor() commands.CommandExecutor {
	return func(args []string) (output.Command, error) {

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

		var filtered []upcloud.Network
		for _, n := range networks.Networks {
			if s.all {
				filtered = append(filtered, n)
				continue
			}

			if s.public && n.Type == upcloud.NetworkTypePublic {
				filtered = append(filtered, n)
			}
			if s.utility && n.Type == upcloud.NetworkTypeUtility {
				filtered = append(filtered, n)
			}
			if !s.public && !s.utility && n.Type == upcloud.NetworkTypePrivate {
				filtered = append(filtered, n)
			}
		}

		var rows []output.TableRow
		for _, n := range filtered {
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
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.zone, "zone", "", "Show networks from a specific zone.")
	flags.BoolVar(&s.all, "all", false, "Show all networks.")
	flags.BoolVar(&s.public, "public", false, "Show public networks instead of private networks.")
	flags.BoolVar(&s.utility, "utility", false, "Show utility networks instead of private networks.")
	//	flags.BoolVar(&s.private, "private", true, "Show private networks (default).")
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}
