package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// PlansCommand creates the "loadbalancer plans" command
func PlansCommand() commands.Command {
	return &plansCommand{
		BaseCommand: commands.New("plans", "List available load balancer plans", "upctl loadbalancer plans"),
	}
}

type plansCommand struct {
	*commands.BaseCommand
}

// Execute implements commands.NoArgumentCommand
func (s *plansCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	plans, err := svc.GetLoadBalancerPlans(exec.Context(), &request.GetLoadBalancerPlansRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, plan := range plans {
		rows = append(rows, output.TableRow{
			plan.Name,
			plan.PerServerMaxSessions,
			plan.ServerNumber,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "per_server_max_sessions", Header: "Max sessions per server"},
				{Key: "server_number", Header: "Server count"},
			},
			Rows: rows,
		},
	}, nil
}
