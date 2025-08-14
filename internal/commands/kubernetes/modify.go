package kubernetes

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ModifyCommand creates the "kubernetes modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modify an existing cluster",
			"upctl kubernetes modify 00bb4617-c592-4b32-b869-35a60b323b18 --plan 1xCPU-1GB",
		),
	}
}

type modifyCommand struct {
	*commands.BaseCommand
	resolver.CachingKubernetes
	completion.Kubernetes

	controlPlaneIPFilter []string
	labels               []string
	clearLabels          config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (c *modifyCommand) InitCommand() {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"uks"})

	fs := &pflag.FlagSet{}
	fs.StringArrayVar(
		&c.controlPlaneIPFilter,
		"kubernetes-api-allow-ip",
		[]string{},
		"Allow cluster's Kubernetes API to be accessed from an IP address or a network CIDR, multiple can be declared.",
	)
	fs.StringArrayVar(&c.labels, "label", nil, "Labels to describe the cluster in `key=value` format, multiple can be declared.")
	config.AddToggleFlag(fs, &c.clearLabels, "clear-labels", false, "Clear all labels from to given cluster.")

	c.AddFlags(fs)
	c.Cobra().MarkFlagsMutuallyExclusive("label", "clear-labels")
	for _, flag := range []string{"kubernetes-api-allow-ip", "label"} {
		commands.Must(c.Cobra().RegisterFlagCompletionFunc(flag, cobra.NoFileCompletions))
	}
}

// Execute implements commands.MultipleArgumentCommand
func (c *modifyCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"uks"}, "k8s")

	msg := fmt.Sprintf("Modifying Kubernetes cluster %v", arg)
	exec.PushProgressStarted(msg)

	req := request.ModifyKubernetesClusterRequest{
		ClusterUUID: arg,
		Cluster:     request.ModifyKubernetesCluster{},
	}

	if len(c.controlPlaneIPFilter) > 0 {
		req.Cluster.ControlPlaneIPFilter = &c.controlPlaneIPFilter
	}

	if c.clearLabels.Value() {
		req.Cluster.Labels = &[]upcloud.Label{}
	}

	if len(c.labels) > 0 {
		labelSlice, err := labels.StringsToSliceOfLabels(c.labels)
		if err != nil {
			return nil, err
		}

		req.Cluster.Labels = &labelSlice
	}

	res, err := exec.All().ModifyKubernetesCluster(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
