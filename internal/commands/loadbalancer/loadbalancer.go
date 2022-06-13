package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
)

// BaseLoadBalancerCommand creates the base "loadbalancer" command
func BaseLoadBalancerCommand() commands.Command {
	return &loadbalancerCommand{
		commands.New("loadbalancer", "Manage load balancers"),
	}
}

type loadbalancerCommand struct {
	*commands.BaseCommand
}
