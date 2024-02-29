package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the "loadbalancer list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current load balancers", "upctl loadbalancer list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	loadbalancers, err := svc.GetLoadBalancers(exec.Context(), &request.GetLoadBalancersRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, lb := range loadbalancers {
		rows = append(rows, output.TableRow{
			lb.UUID,
			lb.Name,
			lb.Plan,
			lb.Zone,
			lb.OperationalState,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: loadbalancers,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "plan", Header: "Plan"},
				{Key: "zone", Header: "Zone"},
				{Key: "operational_state", Header: "State", Format: format.LoadBalancerState},
			},
			Rows: rows,
		},
	}, nil
}
