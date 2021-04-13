package router

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/pflag"
)

// ListCommand creates the "router list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List routers", ""),
	}
}

type listCommand struct {
	*commands.BaseCommand
	allRouters     bool
	normalRouters  bool
	serviceRouters bool
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	//	s.header = table.Row{"UUID", "Name", "Type"}
	flags := &pflag.FlagSet{}
	flags.BoolVar(&s.allRouters, "all", false, "Show all routers.")
	flags.BoolVar(&s.normalRouters, "normal", true, "Show normal routers (default).")
	flags.BoolVar(&s.serviceRouters, "service", false, "Show service routers.")
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

	if s.serviceRouters {
		s.normalRouters = false
	}
	var filtered []upcloud.Router
	if s.allRouters {
		filtered = routers.Routers
	} else {
		for _, r := range routers.Routers {
			if s.normalRouters && r.Type == "normal" {
				filtered = append(filtered, r)
			}
			if s.serviceRouters && r.Type == "service" {
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
			{Header: "UUID", Key: "uuid", Color: ui.DefaultUUUIDColours},
			{Header: "Name", Key: "name"},
			{Header: "Type", Key: "type"},
		},
		Rows: rows,
	}, nil
}
