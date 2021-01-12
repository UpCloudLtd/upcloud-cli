package server

import (
	"fmt"
	"io"
	"sort"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func PlanListCommand() commands.Command {
	return &planListCommand{
		BaseCommand: commands.New("plans", "List Server Plans"),
	}
}

type planListCommand struct {
	*commands.BaseCommand
	header         table.Row
	columnKeys     []string
	visibleColumns []string
}

func (s *planListCommand) InitCommand() {
	s.header = table.Row{"Name", "Cores", "Memory (MiB)", "Storage (GiB)", "Storage tier", "Traffic out / month (MiB)"}
	s.columnKeys = []string{"name", "cores", "memory", "storage", "storage_tier", "traffic"}
	s.visibleColumns = []string{"name", "cores", "memory", "storage", "storage_tier", "traffic"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *planListCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		service := upapi.Service(s.Config())
		plans, err := service.GetPlans()
		if err != nil {
			return nil, err
		}
		return plans, nil
	}
}

func (s *planListCommand) HandleOutput(writer io.Writer, out interface{}) error {
	plans := out.(*upcloud.Plans)
	t := ui.NewDataTable(s.columnKeys...)
	t.SetHeader(s.header)
	t.OverrideColumnKeys(s.visibleColumns...)

	sort.SliceStable(plans.Plans, func(i, j int) bool {
		return plans.Plans[i].MemoryAmount < plans.Plans[j].MemoryAmount == true
	})
	sort.SliceStable(plans.Plans, func(i, j int) bool {
		return plans.Plans[i].CoreNumber < plans.Plans[j].CoreNumber == true
	})

	for _, plan := range plans.Plans {
		t.AppendRow(table.Row{
			plan.Name,
			fmt.Sprintf("%d", plan.CoreNumber),
			fmt.Sprintf("%d", plan.MemoryAmount),
			fmt.Sprintf("%d", plan.StorageSize),
			plan.StorageTier,
			fmt.Sprintf("%d", plan.PublicTrafficOut),
		})
	}

	return t.Paginate(writer)
}
