package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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
	loadbalancers, err := svc.GetLoadBalancers(&request.GetLoadBalancersRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, lb := range loadbalancers {
		coloredState := commands.LoadBalancerOperationalStateColour(lb.OperationalState).Sprint(lb.OperationalState)

		rows = append(rows, output.TableRow{
			lb.UUID,
			lb.Name,
			lb.Plan,
			lb.Zone,
			coloredState,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
			{Key: "name", Header: "Name"},
			{Key: "plan", Header: "Plan"},
			{Key: "zone", Header: "Zone"},
			{Key: "state", Header: "State"},
		},
		Rows: rows,
	}, nil
}
