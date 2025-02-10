package nodegroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ShowCommand creates the "kubernetes nodegroup show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show node group details",
			"upctl kubernetes nodegroup show 55199a44-4751-4e27-9394-7c7661910be3 --name default",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	name string
	resolver.CachingKubernetes
	completion.Kubernetes
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.name, "name", "", "Node group name")
	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("name", cobra.NoFileCompletions))
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *showCommand) ExecuteSingleArgument(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	nodeGroup, err := svc.GetKubernetesNodeGroup(exec.Context(), &request.GetKubernetesNodeGroupRequest{ClusterUUID: uuid, Name: s.name})
	if err != nil {
		return nil, err
	}

	taintColumns := []output.TableColumn{
		{Key: "key", Header: "Key"},
		{Key: "value", Header: "Value"},
		{Key: "effect", Header: "Effect"},
	}

	taintRows := []output.TableRow{}
	for _, taint := range nodeGroup.Taints {
		taintRows = append(taintRows, output.TableRow{
			taint.Key,
			taint.Value,
			taint.Effect,
		})
	}

	kubeletArgColumns := taintColumns[0:2]

	kubeletArgRows := []output.TableRow{}
	for _, kubeletArg := range nodeGroup.KubeletArgs {
		kubeletArgRows = append(kubeletArgRows, output.TableRow{
			kubeletArg.Key,
			kubeletArg.Value,
		})
	}

	nodeRows := []output.TableRow{}
	for _, node := range nodeGroup.Nodes {
		nodeRows = append(nodeRows, output.TableRow{
			node.UUID,
			node.Name,
			node.State,
		})
	}

	nodeColumns := []output.TableColumn{
		{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
		{Key: "name", Header: "Name"},
		{Key: "state", Header: "State", Format: format.KubernetesNodeState},
	}

	var storageName string
	if storage, err := svc.GetStorageDetails(exec.Context(), &request.GetStorageDetailsRequest{UUID: nodeGroup.Storage}); err != nil {
		storageName = ""
	} else {
		storageName = storage.Title
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: nodeGroup,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: []output.DetailSection{
						{
							Title: "Overview",
							Rows: []output.DetailRow{
								{Title: "Name:", Value: nodeGroup.Name},
								{Title: "Count:", Value: nodeGroup.Count},
								{Title: "Plan:", Value: nodeGroup.Plan},
								{Title: "State:", Value: nodeGroup.State, Format: format.KubernetesNodeGroupState},
								{Title: "Storage UUID:", Value: nodeGroup.Storage, Colour: ui.DefaultUUUIDColours},
								{Title: "Storage name:", Value: storageName, Format: format.PossiblyUnknownString},
								{Title: "Anti-affinity:", Value: nodeGroup.AntiAffinity, Format: format.Boolean},
								{Title: "Utility network access:", Value: nodeGroup.UtilityNetworkAccess, Format: format.Boolean},
							},
						},
					},
				},
			},
			labels.GetLabelsSectionWithResourceType(nodeGroup.Labels, "node group"),
			output.CombinedSection{
				Key:   "taints",
				Title: "Taints:",
				Contents: output.Table{
					Columns:      taintColumns,
					Rows:         taintRows,
					EmptyMessage: "No taints defined for this node group.",
				},
			},
			output.CombinedSection{
				Key:   "kubelet_args",
				Title: "Kubelet arguments:",
				Contents: output.Table{
					Columns:      kubeletArgColumns,
					Rows:         kubeletArgRows,
					EmptyMessage: "No kubelet arguments defined for this node group.",
				},
			},
			output.CombinedSection{
				Key:   "nodes",
				Title: "Nodes:",
				Contents: output.Table{
					Columns:      nodeColumns,
					Rows:         nodeRows,
					EmptyMessage: "No nodes found for this node group.",
				},
			},
		},
	}, nil
}
