package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud/request"
)

// PlansCommand creates the "loadbalancer plans" command
func PlansCommand() commands.Command {
	return &plansCommand{
		BaseCommand: commands.New("plans", "List available cluster plans", "upctl kubernetes plans"),
	}
}

type plansCommand struct {
	*commands.BaseCommand
}

// Execute implements commands.NoArgumentCommand
func (s *plansCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	plans, err := svc.GetKubernetesPlans(exec.Context(), &request.GetKubernetesPlansRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, plan := range plans {
		rows = append(rows, output.TableRow{
			plan.Name,
			plan.ServerNumber,
			plan.MaxNodes,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "server_number", Header: "Control nodes"},
				{Key: "max_nodes", Header: "Max worker nodes"},
			},
			Rows: rows,
		},
	}, nil
}
