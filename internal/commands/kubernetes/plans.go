package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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

func (s *plansCommand) InitCommand() {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"uks"})
}

// Execute implements commands.NoArgumentCommand
func (s *plansCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"uks"}, "k8s")

	svc := exec.All()
	plans, err := svc.GetKubernetesPlans(exec.Context(), &request.GetKubernetesPlansRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, plan := range plans {
		if !plan.Deprecated {
			rows = append(rows, output.TableRow{
				plan.Name,
				plan.MaxNodes,
			})
		}
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "max_nodes", Header: "Max worker nodes"},
			},
			Rows: rows,
		},
	}, nil
}
