package loadbalancer

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	err := svc.DeleteLoadBalancer(&request.DeleteLoadBalancerRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}
	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	return output.None{}, nil
}
