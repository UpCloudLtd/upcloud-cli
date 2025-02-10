package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// PlansCommand creates the "loadbalancer plans" command
func PlansCommand() commands.Command {
	return &plansCommand{
		BaseCommand: commands.New("plans", "List available load balancer plans", "upctl load-balancer plans"),
	}
}

type plansCommand struct {
	*commands.BaseCommand
}

func (s *plansCommand) InitCommand() {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"loadbalancer"})
}

// Execute implements commands.NoArgumentCommand
func (s *plansCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"loadbalancer"}, "load-balancer")

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
