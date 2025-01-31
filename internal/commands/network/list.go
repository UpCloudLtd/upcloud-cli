package network

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// ListCommand creates the "network list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List networks, by default private networks only",
			"upctl network list",
			"upctl network list --zone pl-waw1",
			"upctl network list --zone pl-waw1 --public",
			"upctl network list --all",
			"upctl network list --zone pl-waw1 --all",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	zone    string
	all     config.OptionalBoolean
	public  config.OptionalBoolean
	utility config.OptionalBoolean
}

func (s *listCommand) MaximumExecutions() int {
	return 1
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.zone, "zone", "", "Show networks from a specific zone.")
	config.AddToggleFlag(flags, &s.all, "all", false, "Show all networks.")
	config.AddToggleFlag(flags, &s.public, "public", false, "Show public networks instead of private networks.")
	config.AddToggleFlag(flags, &s.utility, "utility", false, "Show utility networks instead of private networks.")
	// TODO: reimplmement
	// s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *listCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Network()
	var networks *upcloud.Networks
	var err error
	if s.zone != "" {
		networks, err = svc.GetNetworksInZone(exec.Context(), &request.GetNetworksInZoneRequest{Zone: s.zone})
	} else {
		networks, err = svc.GetNetworks(exec.Context())
	}
	if err != nil {
		return nil, err
	}

	var filtered []upcloud.Network
	for _, n := range networks.Networks {
		if s.all.Value() {
			filtered = append(filtered, n)
			continue
		}

		if s.public.Value() && n.Type == upcloud.NetworkTypePublic {
			filtered = append(filtered, n)
		}
		if s.utility.Value() && n.Type == upcloud.NetworkTypeUtility {
			filtered = append(filtered, n)
		}
		if !s.public.Value() && !s.utility.Value() && n.Type == upcloud.NetworkTypePrivate {
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

	return output.MarshaledWithHumanOutput{
		Value: upcloud.Networks{Networks: filtered},
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "router", Header: "Router", Colour: ui.DefaultUUUIDColours},
				{Key: "type", Header: "Type"},
				{Key: "zone", Header: "Zone"},
			},
			Rows: rows,
		},
	}, nil
}
