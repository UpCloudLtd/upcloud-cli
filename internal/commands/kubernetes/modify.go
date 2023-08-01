package kubernetes

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/pflag"
)

// ModifyCommand creates the "kubernetes modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modifiy an existing cluster",
			"upctl cluster modify 00bb4617-c592-4b32-b869-35a60b323b18 --plan 1xCPU-1GB",
		),
	}
}

type modifyCommand struct {
	*commands.BaseCommand
	resolver.CachingKubernetes
	completion.Kubernetes
	params request.ModifyKubernetesClusterRequest
}

// InitCommand implements Command.InitCommand
func (c *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringArrayVar(
		&c.params.Cluster.ControlPlaneIPFilter,
		"control-plane-allow-ip",
		[]string{},
		"Allow cluster control-plane to be accessed from an IP address or a network CIDR, multiple can be declared..",
	)

	c.AddFlags(fs)
}

// Execute implements commands.MultipleArgumentCommand
func (c *modifyCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Modifying Kubernetes cluster %v", arg)
	exec.PushProgressStarted(msg)

	req := c.params
	req.ClusterUUID = arg

	res, err := exec.All().ModifyKubernetesCluster(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
