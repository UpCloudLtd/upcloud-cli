package nodegroup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type deleteCommand struct {
	*commands.BaseCommand
	name string
	completion.Kubernetes
	resolver.CachingKubernetes
}

// DeleteCommand creates the "kubernetes nodegroup delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete the node group from the cluster.",
			"upctl kubernetes nodegroup delete 55199a44-4751-4e27-9394-7c7661910be3 --name secondary-node-group",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.name, "name", "", "Node group name")
	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("name", cobra.NoFileCompletions))
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *deleteCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Deleting node group %s from cluster %v", s.name, arg)
	exec.PushProgressStarted(msg)

	err := exec.All().DeleteKubernetesNodeGroup(exec.Context(), &request.DeleteKubernetesNodeGroupRequest{
		ClusterUUID: arg, Name: s.name,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
