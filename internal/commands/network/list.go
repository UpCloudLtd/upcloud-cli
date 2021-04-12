package network

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

// ListCommand creates the "network list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List networks, by default private networks only"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	zone    string
	all     bool
	public  bool
	utility bool
}

func (s *listCommand) MaximumExecutions() int {
	return 1
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.zone, "zone", "", "Show networks from a specific zone.")
	flags.BoolVar(&s.all, "all", false, "Show all networks.")
	flags.BoolVar(&s.public, "public", false, "Show public networks instead of private networks.")
	flags.BoolVar(&s.utility, "utility", false, "Show utility networks instead of private networks.")
	//	flags.BoolVar(&s.private, "private", true, "Show private networks (default).")
	// TODO: reimplmement
	// s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

// Execute implements command.Command
func (s *listCommand) Execute(exec commands.Executor, _ string) (output.Output, error) {
	svc := exec.Network()
	var networks *upcloud.Networks
	var err error
	if s.zone != "" {
		networks, err = svc.GetNetworksInZone(&request.GetNetworksInZoneRequest{Zone: s.zone})
	} else {
		networks, err = svc.GetNetworks()
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
			n.UUID,
			n.Name,
			n.Router,
			n.Type,
			n.Zone,
		})
	}
	return output.Table{
		Columns: []output.TableColumn{
			{Key: "uuid", Header: "UUID", Color: ui.DefaultUUUIDColours},
			{Key: "name", Header: "Name"},
			{Key: "router", Header: "Router"},
			{Key: "type", Header: "Type"},
			{Key: "zone", Header: "Zone"},
		},
		Rows: rows,
	}, nil
}
