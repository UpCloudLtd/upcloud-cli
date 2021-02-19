package server

import (
	"fmt"
	"io"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

// ListCommand creates the "server list" command
func ListCommand(service service.Server) commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current servers"),
		service:     service,
	}
}

type listCommand struct {
	*commands.BaseCommand
	service        service.Server
	header         table.Row
	columnKeys     []string
	visibleColumns []string
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	s.header = table.Row{"UUID", "Hostname", "Plan", "Zone", "State", "Tags", "Title", "Licence"}
	s.columnKeys = []string{"uuid", "hostname", "plan", "zone", "state", "tags", "title", "licence"}
	s.visibleColumns = []string{"uuid", "hostname", "plan", "zone", "state"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		svc := upapi.Service(s.Config())
		servers, err := svc.GetServers()
		if err != nil {
			return nil, err
		}
		return servers, nil
	}
}

// HandleOutput implements Command.HandleOutput
func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
	servers := out.(*upcloud.Servers)
	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	t.SetColumnConfig("state", table.ColumnConfig{Transformer: func(val interface{}) string {
		return commands.StateColour(val.(string)).Sprint(val)
	}})

	for _, server := range servers.Servers {
		plan := server.Plan
		if plan == "custom" {
			memory := server.MemoryAmount / 1024
			plan = fmt.Sprintf("Custom (%dxCPU, %dGB)", server.CoreNumber, memory)
		}
		t.AppendRow(table.Row{
			server.UUID,
			server.Hostname,
			plan,
			server.Zone,
			server.State,
			strings.Join(server.Tags, ","),
			server.Title,
			server.License})
	}

	return t.Paginate(writer)
}
