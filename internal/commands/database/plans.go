package database

import (
	"sort"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// PlansCommand creates the "database plans" command
func PlansCommand() commands.Command {
	return &plansCommand{
		BaseCommand: commands.New("plans", "List available plans for given database type", "upctl database plans pg", "upctl database plans mysql"),
	}
}

type plansCommand struct {
	*commands.BaseCommand
	completion.DatabaseType
}

// Execute implements commands.MultipleArgumentCommand
func (s *plansCommand) Execute(exec commands.Executor, serviceType string) (output.Output, error) {
	svc := exec.All()
	dbType, err := svc.GetManagedDatabaseServiceType(exec.Context(), &request.GetManagedDatabaseServiceTypeRequest{Type: serviceType})
	if err != nil {
		return nil, err
	}

	plans := dbType.ServicePlans
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].NodeCount != plans[j].NodeCount {
			return plans[i].NodeCount < plans[j].NodeCount
		}

		if plans[i].CoreNumber != plans[j].CoreNumber {
			return plans[i].CoreNumber < plans[j].CoreNumber
		}

		if plans[i].MemoryAmount != plans[j].MemoryAmount {
			return plans[i].MemoryAmount < plans[j].MemoryAmount
		}

		return plans[i].StorageSize < plans[j].StorageSize
	})

	rows := []output.TableRow{}
	for _, plan := range plans {
		rows = append(rows, output.TableRow{
			plan.Plan,
			plan.NodeCount,
			plan.CoreNumber,
			plan.MemoryAmount / 1024,
			plan.StorageSize / 1024,
			plan.BackupConfig.Interval,
			plan.BackupConfig.MaxCount,
			plan.BackupConfig.RecoveryMode,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "plan", Header: "Plan"},
				{Key: "nodes", Header: "Nodes"},
				{Key: "cores", Header: "Cores"},
				{Key: "memory", Header: "Memory (GB)"},
				{Key: "storage", Header: "Storage (GB)"},
				{Key: "bu_interval", Header: "Backup interval"},
				{Key: "bu_max_count", Header: "Max backup count"},
			},
			Rows: rows,
		},
	}, nil
}
