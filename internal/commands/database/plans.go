package database

import (
	"fmt"
	"sort"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
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
	showLegacy config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (s *plansCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &s.showLegacy, "show-legacy", false, "List also legacy plans.")
	s.AddFlags(flags)
}

// Execute implements commands.MultipleArgumentCommand
func (s *plansCommand) Execute(exec commands.Executor, serviceType string) (output.Output, error) {
	svc := exec.All()

	plans, err := svc.GetManagedDatabasePlans(exec.Context(), &request.GetManagedDatabasePlansRequest{})
	if err != nil {
		return nil, err
	}

	computeShapes := make([]upcloud.ManagedDatabasePlanComputeShape, 0)

	for _, plan := range plans {
		if plan.Type == serviceType {
			computeShapes = plan.ComputeShapes
		}
	}

	var legacyPlans []upcloud.ManagedDatabaseServicePlan
	if s.showLegacy.Value() == true {
		dbType, err := svc.GetManagedDatabaseServiceType(exec.Context(), &request.GetManagedDatabaseServiceTypeRequest{Type: serviceType})
		if err != nil {
			return nil, err
		}

		legacyPlans = dbType.ServicePlans
		sort.Slice(legacyPlans, func(i, j int) bool {
			if legacyPlans[i].NodeCount != legacyPlans[j].NodeCount {
				return legacyPlans[i].NodeCount < legacyPlans[j].NodeCount
			}

			if legacyPlans[i].CoreNumber != legacyPlans[j].CoreNumber {
				return legacyPlans[i].CoreNumber < legacyPlans[j].CoreNumber
			}

			if legacyPlans[i].MemoryAmount != legacyPlans[j].MemoryAmount {
				return legacyPlans[i].MemoryAmount < legacyPlans[j].MemoryAmount
			}

			return legacyPlans[i].StorageSize < legacyPlans[j].StorageSize
		})

		for _, plan := range legacyPlans {
			computeShapes = append(computeShapes, upcloud.ManagedDatabasePlanComputeShape{
				Compute:                 plan.Plan,
				Family:                  "legacy",
				CPU:                     plan.CoreNumber,
				MemoryGB:                plan.MemoryAmount,
				DynamicStorageSupported: false,
				NodeCounts:              []int{plan.NodeCount},
				Backups:                 []string{fmt.Sprintf("%d PITR days", plan.BackupConfig.MaxCount)},
				Storage: upcloud.ManagedDatabasePlanStorageInfo{
					StepGiB:              0,
					DynamicMaxMultiplier: 0,
					TotalCapGiB:          0,
					Options: []upcloud.ManagedDatabasePlanStorageOption{
						{BaseGiB: plan.StorageSize, MaxGiB: plan.StorageSize},
					},
				},
			})
		}

	}

	rows := []output.TableRow{}
	for _, computeShape := range computeShapes {
		nodes := make([]string, 0)
		for _, nodeCount := range computeShape.NodeCounts {
			nodes = append(nodes, fmt.Sprintf("%d", nodeCount))
		}

		storageRanges := make([]string, 0)
		for _, option := range computeShape.Storage.Options {
			storageRange := fmt.Sprintf("%d", option.BaseGiB)
			if option.MaxGiB != option.BaseGiB {
				storageRange = fmt.Sprintf("%s-%d", storageRange, option.MaxGiB)
			}
			storageRanges = append(storageRanges, storageRange)
		}

		rows = append(rows, output.TableRow{
			computeShape.Compute,
			computeShape.Family,
			computeShape.CPU,
			computeShape.MemoryGB,
			strings.Join(nodes, ", "),
			strings.Join(storageRanges, ", "),
			strings.Join(computeShape.Backups, ", "),
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: computeShapes,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "compute_shape", Header: "Compute shape"},
				{Key: "family", Header: "Family"},
				{Key: "cores", Header: "Cores"},
				{Key: "memory", Header: "Memory (GB)"},
				{Key: "nodes", Header: "Nodes"},
				{Key: "storage_ranges", Header: "Storage ranges (GB)"},
				{Key: "backups", Header: "Backups"},
			},
			Rows: rows,
		},
	}, nil
}
