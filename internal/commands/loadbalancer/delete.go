package loadbalancer

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "loadbalancer delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a load balancer",
			"upctl load-balancer delete 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl load-balancer delete my-load-balancer",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingLoadBalancer
	completion.LoadBalancer
}

func (s *deleteCommand) InitCommand() {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"loadbalancer"})
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"loadbalancer"}, "load-balancer")

	svc := exec.All()
	msg := fmt.Sprintf("Deleting load balancer %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteLoadBalancer(exec.Context(), &request.DeleteLoadBalancerRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, fmt.Sprintf("%s: failed", msg), err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
