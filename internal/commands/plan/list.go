package plan

import (
	"fmt"
	"sort"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List Server Plans"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	header         table.Row
	columnKeys     []string
	visibleColumns []string
}

func (s *listCommand) InitCommand() {
	s.header = table.Row{"Name", "Cores", "Memory (MiB)", "Storage (GiB)", "Storage tier", "Traffic out / month (MiB)"}
	s.columnKeys = []string{"name", "cores", "memory", "storage", "storage_tier", "traffic"}
	s.visibleColumns = []string{"name", "cores", "memory", "storage", "storage_tier", "traffic"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		service := upapi.Service(s.Config())
		plans, err := service.GetPlans()
		if err != nil {
			return nil, err
		}
		return plans, nil
	}
}

func (s *listCommand) HandleOutput(out interface{}) error {
	plans := out.(*upcloud.Plans)
	fmt.Println()
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
	fmt.Println(t.Render())
	fmt.Println()
	return nil
}
