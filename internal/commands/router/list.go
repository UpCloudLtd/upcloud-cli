package router

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/spf13/pflag"
)

// ListCommand creates the "router list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List routers",
			"upctl router list",
			"upctl router list --all",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	allRouters     config.OptionalBoolean
	normalRouters  config.OptionalBoolean
	serviceRouters config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	//	s.header = table.Row{"UUID", "Name", "Type"}
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &s.allRouters, "all", false, "Show all routers.")
	config.AddToggleFlag(flags, &s.normalRouters, "normal", true, "Show normal routers.")
	config.AddToggleFlag(flags, &s.serviceRouters, "service", false, "Show service routers.")
	// TODO: reimplement
	// 	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	routers, err := exec.Network().GetRouters()
	if err != nil {
		return nil, err
	}

	if s.serviceRouters.Value() {
		s.normalRouters = config.False
	}
	var filtered []upcloud.Router
	if s.allRouters.Value() {
		filtered = routers.Routers
	} else {
		for _, r := range routers.Routers {
			if s.normalRouters.Value() && r.Type == "normal" {
				filtered = append(filtered, r)
			}
			if s.serviceRouters.Value() && r.Type == "service" {
				filtered = append(filtered, r)
			}
		}
	}
	rows := make([]output.TableRow, len(filtered))
	for i, router := range filtered {
		rows[i] = output.TableRow{router.UUID, router.Name, router.Type}
	}
	return output.Table{
		Columns: []output.TableColumn{
			{Header: "UUID", Key: "uuid", Colour: ui.DefaultUUUIDColours},
			{Header: "Name", Key: "name"},
			{Header: "Type", Key: "type"},
		},
		Rows: rows,
	}, nil
}
