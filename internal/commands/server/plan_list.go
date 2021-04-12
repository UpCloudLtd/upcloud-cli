package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
)

// PlanListCommand creates the "server plans" command
func PlanListCommand() commands.Command {
	return &planListCommand{
		BaseCommand: commands.New("plans", "List server plans"),
	}
}

type planListCommand struct {
	*commands.BaseCommand
}

func (s *planListCommand) Execute(exec commands.Executor, _ string) (output.Output, error) {
	svc := exec.Plan()
	plans, err := svc.GetPlans()
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, p := range plans.Plans {
		rows = append(rows, output.TableRow{
			p.Name,
			p.CoreNumber,
			p.MemoryAmount,
			p.StorageSize,
			p.StorageTier,
			p.PublicTrafficOut,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "name", Header: "Name"},
			{Key: "cores", Header: "Cores"},
			{Key: "memory", Header: "Memory"},
			{Key: "storage", Header: "Storage size"},
			{Key: "storage_tier", Header: "Storage tier"},
			{Key: "egress_transfer", Header: "Transfer out (GiB/month)"},
		},
		Rows: rows,
	}, nil
}
