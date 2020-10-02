package plan

import (
	"fmt"
	"os"
	"sort"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/table"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func ListCommand() commands.Command {
	return &listCommand{
		Command: commands.New("list", "List Server Plans"),
	}
}

type listCommand struct {
	commands.Command
	headerNames    []string
	columnKeys     []string
	visibleColumns []string
}

func (s *listCommand) InitCommand() {
	s.headerNames = []string{"Cores", "Memory (MiB)", "Storage (GiB)", "Storage tier", "Traffic out per month (MiB)"}
	s.columnKeys = []string{"cores", "memory", "storage", "storage_tier", "traffic"}
	s.visibleColumns = []string{"cores", "memory", "storage", "storage_tier", "traffic"}
	flags := &pflag.FlagSet{}
	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	s.AddFlags(flags)
}

func (s *listCommand) MakeExecuteCommand() func(args []string) error {
	return func(args []string) error {
		service := upapi.Service(s.Config())
		plans, err := service.GetPlans()
		if err != nil {
			return err
		}
		s.HandleOutput(plans)
		return nil
	}
}

func (s *listCommand) HandleOutput(out interface{}) error {
	if s.Config().GetString("output") != "human" {
		return s.Command.HandleOutput(out)
	}
	plans := out.(*upcloud.Plans)
	fmt.Println()
	t, err := table.New(os.Stdout, s.headerNames...)
	if err != nil {
		return err
	}
	t.SetColumnKeys(s.columnKeys...)
	t.SetVisibleColumns(s.visibleColumns...)

	sort.SliceStable(plans.Plans, func(i, j int) bool {
		return plans.Plans[i].MemoryAmount < plans.Plans[j].MemoryAmount == true
	})
	sort.SliceStable(plans.Plans, func(i, j int) bool {
		return plans.Plans[i].CoreNumber < plans.Plans[j].CoreNumber == true
	})

	for _, plan := range plans.Plans {
		t.NextRow()
		for _, v := range []string{
			plan.Name,
			fmt.Sprintf("%d", plan.CoreNumber),
			fmt.Sprintf("%d", plan.MemoryAmount),
			fmt.Sprintf("%d", plan.StorageSize),
			plan.StorageTier,
			fmt.Sprintf("%d", plan.PublicTrafficOut),
		} {
			t.AddColumn(v, nil)
		}
	}
	t.Render()
	fmt.Println()
	return nil
}
