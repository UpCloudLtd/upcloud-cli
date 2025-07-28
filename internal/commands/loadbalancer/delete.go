package loadbalancer

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
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

	wait config.OptionalBoolean
}

func (s *deleteCommand) InitCommand() {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"loadbalancer"})

	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &s.wait, "wait", false, "Wait until the Kubernetes cluster has been deleted before returning.")
	s.AddFlags(flags)
}

func Delete(exec commands.Executor, uuid string, wait bool) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting load balancer %v", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteLoadBalancer(exec.Context(), &request.DeleteLoadBalancerRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for load balancer %s to be deleted", uuid))
		err = svc.WaitForLoadBalancerDeletion(exec.Context(), &request.WaitForLoadBalancerDeletionRequest{
			UUID: uuid,
		})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		exec.PushProgressUpdateMessage(msg, msg)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"loadbalancer"}, "load-balancer")

	return Delete(exec, arg, s.wait.Value())
}
