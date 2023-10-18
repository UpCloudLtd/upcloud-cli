package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
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

// InitCommand implements Command.InitCommand
func (lb *loadbalancerCommand) InitCommand() {
	lb.Cobra().Aliases = []string{"lb"}
}
