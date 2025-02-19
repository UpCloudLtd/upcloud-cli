package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseLoadBalancerCommand creates the base "loadbalancer" command
func BaseLoadBalancerCommand() commands.Command {
	// Initialize the BaseCommand properly
	baseCmd := commands.New("load-balancer", "Manage load balancers")
	baseCmd.SetDeprecatedAliases([]string{"loadbalancer"})

	lbc := &loadbalancerCommand{
		BaseCommand: baseCmd,
	}

	return lbc
}

type loadbalancerCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (lb *loadbalancerCommand) InitCommand() {
	lb.Cobra().Aliases = []string{"lb", "loadbalancer"}

	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetDeprecationHelp(lb.Cobra(), lb.DeprecatedAliases())
}
