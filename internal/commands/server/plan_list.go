package server

import (
	"sort"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// PlanListCommand creates the "server plans" command
func PlanListCommand() commands.Command {
	return &planListCommand{
		BaseCommand: commands.New("plans", "List server plans", "upctl server plans"),
	}
}

type planListCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *planListCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	plansObj, err := exec.All().GetPlans(exec.Context())
	if err != nil {
		return nil, err
	}

	plans := plansObj.Plans
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].CoreNumber != plans[j].CoreNumber {
			return plans[i].CoreNumber < plans[j].CoreNumber
		}

		if plans[i].MemoryAmount != plans[j].MemoryAmount {
			return plans[i].MemoryAmount < plans[j].MemoryAmount
		}

		return plans[i].StorageSize < plans[j].StorageSize
	})

	rows := make(map[string][]output.TableRow)
	for _, p := range plans {
		key := planType(p)
		row := output.TableRow{
			p.Name,
			p.CoreNumber,
			p.MemoryAmount,
			p.StorageSize,
			p.StorageTier,
			p.PublicTrafficOut,
		}

		// Add GPU fields only for GPU plans
		if key == "gpu" {
			row = append(row, p.GPUModel, p.GPUAmount)
		}

		rows[key] = append(rows[key], row)
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Combined{
			planSection("general_purpose", "General purpose", rows["general_purpose"]),
			planSection("gpu", "GPU", rows["gpu"]),
			planSection("cloud_native", "Cloud native", rows["cloud_native"]),
			planSection("high_cpu", "High CPU", rows["high_cpu"]),
			planSection("high_memory", "High memory", rows["high_memory"]),
			planSection("developer", "Developer", rows["developer"]),
		},
	}, nil
}

func planType(p upcloud.Plan) string {
	if strings.HasPrefix(p.Name, "DEV-") {
		return "developer"
	}
	if strings.HasPrefix(p.Name, "HICPU-") {
		return "high_cpu"
	}
	if strings.HasPrefix(p.Name, "HIMEM-") {
		return "high_memory"
	}
	if strings.HasPrefix(p.Name, "CLOUDNATIVE-") {
		return "cloud_native"
	}
	if strings.HasPrefix(p.Name, "GPU-") {
		return "gpu"
	}
	return "general_purpose"
}

func planSection(key, title string, rows []output.TableRow) output.CombinedSection {
	columns := []output.TableColumn{
		{Key: "name", Header: "Name"},
		{Key: "cores", Header: "Cores"},
		{Key: "memory", Header: "Memory"},
		{Key: "storage", Header: "Storage size"},
		{Key: "storage_tier", Header: "Storage tier"},
		{Key: "egress_transfer", Header: "Transfer out (GiB/month)"},
	}

	if key == "gpu" {
		columns = append(columns,
			output.TableColumn{Key: "gpu_model", Header: "GPU model"},
			output.TableColumn{Key: "gpu_amount", Header: "GPU amount"},
		)
	}

	return output.CombinedSection{
		Key:   key,
		Title: title,
		Contents: output.Table{
			Columns: columns,
			Rows:    rows,
		},
	}
}
