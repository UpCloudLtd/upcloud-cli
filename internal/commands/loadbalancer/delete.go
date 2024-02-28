package loadbalancer

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud/request"
)

// DeleteCommand creates the "loadbalancer delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a load balancer",
			"upctl loadbalancer delete 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl loadbalancer delete my-load-balancer",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingLoadBalancer
	completion.LoadBalancer
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
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
