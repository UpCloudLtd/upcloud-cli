package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"
	"io"
)

type configurationsCommand struct {
	*commands.BaseCommand
	service service.Server
	header         table.Row
	columnKeys     []string
	visibleColumns []string
}

func ConfigurationsCommand(service service.Server) commands.Command {
	return &configurationsCommand{
		BaseCommand: commands.New("configurations", "Lists available server configurations"),
		service:     service,
	}
}

func (s *configurationsCommand) InitCommand() {
	s.header = table.Row{"Number of cores", "Memory amount (Mb)"}
	s.columnKeys = []string{"cores", "memory"}
	s.visibleColumns = []string{"cores", "memory"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *configurationsCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		configurations, err := s.service.GetServerConfigurations()
		if err != nil {
			return nil, err
		}
		return configurations, nil
	}
}

func (s *configurationsCommand) HandleOutput(writer io.Writer, out interface{}) error {
	configurations := out.(*upcloud.ServerConfigurations)

	t := ui.NewDataTable(s.columnKeys...)
	t.OverrideColumnKeys(s.visibleColumns...)
	t.SetHeader(s.header)

	for _, cfg := range configurations.ServerConfigurations {
		t.AppendRow(table.Row{
			cfg.CoreNumber,
			cfg.MemoryAmount,
		})
	}

	return t.Paginate(writer)
}
