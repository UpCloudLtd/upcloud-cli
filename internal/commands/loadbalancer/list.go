package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/paging"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// ListCommand creates the "loadbalancer list" command
func ListCommand() commands.Command {
	cmd := &listCommand{
		BaseCommand: commands.New("list", "List current load balancers", "upctl load-balancer list"),
	}

	return cmd
}

type listCommand struct {
	*commands.BaseCommand
	paging.PageParameters
}

func (s *listCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	s.ConfigureFlags(fs)
	s.AddFlags(fs)
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"loadbalancer"})
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"loadbalancer"}, "load-balancer")
	svc := exec.All()
	loadbalancers, err := svc.GetLoadBalancers(exec.Context(), &request.GetLoadBalancersRequest{
		Page: s.Page(),
	})
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
