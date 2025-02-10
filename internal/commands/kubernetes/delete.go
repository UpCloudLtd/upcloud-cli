package kubernetes

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "kubernetes delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a Kubernetes cluster",
			"upctl kubernetes delete 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl kubernetes delete my-kubernetes-cluster",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingKubernetes
	completion.Kubernetes
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	// Deprecating k8s
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"k8s"})
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating k8s
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"k8s"}, "uks")

	svc := exec.All()
	msg := fmt.Sprintf("Deleting Kubernetes cluster %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteKubernetesCluster(exec.Context(), &request.DeleteKubernetesClusterRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
