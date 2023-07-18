package nodegroup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/pflag"
)

type scaleCommand struct {
	*commands.BaseCommand
	name  string
	count int
	completion.Kubernetes
	resolver.CachingKubernetes
}

// ScaleCommand creates the "kubernetes nodegroup scale" command
func ScaleCommand() commands.Command {
	return &scaleCommand{
		BaseCommand: commands.New(
			"scale",
			"Scale the number of nodes in the node group.",
			"upctl kubernetes nodegroup scale 55199a44-4751-4e27-9394-7c7661910be3 --name secondary-node-group --count 3",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *scaleCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.name, "name", "", "Node group name")
	flagSet.IntVar(&s.count, "count", 0, "Node count")
	s.AddFlags(flagSet)

	_ = s.Cobra().MarkFlagRequired("name")
	_ = s.Cobra().MarkFlagRequired("count")
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *scaleCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Scaling node group %s of cluster %v", s.name, arg)
	exec.PushProgressStarted(msg)

	res, err := exec.All().ModifyKubernetesNodeGroup(exec.Context(), &request.ModifyKubernetesNodeGroupRequest{
		ClusterUUID: arg,
		Name:        s.name,
		NodeGroup:   request.ModifyKubernetesNodeGroup{Count: s.count},
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
