package router

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

// ListCommand creates the "router list" command
func ListCommand(service service.Network) commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List routers"),
		service:     service,
	}
}

type listCommand struct {
	*commands.BaseCommand
	service        service.Network
	header         table.Row
	columnKeys     []string
	visibleColumns []string
	allRouters     bool
	normalRouters  bool
	serviceRouters bool
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	s.header = table.Row{"UUID", "Name", "Type"}
	s.columnKeys = []string{"uuid", "name", "type"}
	s.visibleColumns = []string{"uuid", "name", "type"}
	flags := &pflag.FlagSet{}
	flags.BoolVar(&s.allRouters, "all", false, "Show all routers.")
	flags.BoolVar(&s.normalRouters, "normal", true, "Show normal routers (default).")
	flags.BoolVar(&s.serviceRouters, "service", false, "Show service routers.")
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		routers, err := s.service.GetRouters()
		if err != nil {
			return nil, err
		}

		if s.serviceRouters {
			s.normalRouters = false
		}

		var filtered []upcloud.Router
		for _, r := range routers.Routers {
			if s.allRouters {
				filtered = append(filtered, r)
				continue
			}

			if s.normalRouters && r.Type == "normal" {
				filtered = append(filtered, r)
			}
			if s.serviceRouters && r.Type == "service" {
				filtered = append(filtered, r)
			}
		}

		return &upcloud.Routers{Routers: filtered}, nil
	}
}

// HandleOutput implements Command.HandleOutput
func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
	routers := out.(*upcloud.Routers)

	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	for _, r := range routers.Routers {
		t.Append(table.Row{
			ui.DefaultUUUIDColours.Sprint(r.UUID),
			r.Name,
			r.Type,
		})
	}

	_, _ = fmt.Fprintln(writer, t.Render())
	return nil
}
