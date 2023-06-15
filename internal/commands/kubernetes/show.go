package kubernetes

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// ShowCommand creates the "kubernetes show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show Kubernetes cluster details",
			"upctl kubernetes show 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl kubernetes show my-cluster",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingKubernetes
	completion.Kubernetes
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	cluster, err := svc.GetKubernetesCluster(exec.Context(), &request.GetKubernetesClusterRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	var networkName string
	if network, err := svc.GetNetworkDetails(exec.Context(), &request.GetNetworkDetailsRequest{UUID: cluster.Network}); err != nil {
		networkName = ""
	} else {
		networkName = network.Name
	}

	nodeGroups := output.Combined{}

	for i, nodeGroup := range cluster.NodeGroups {
		kubeletArgs := strings.Builder{}
		for _, v := range nodeGroup.KubeletArgs {
			if kubeletArgs.Len() > 0 {
				kubeletArgs.WriteString("\n")
			}
			_, _ = kubeletArgs.WriteString(fmt.Sprintf("%s=%s", v.Key, v.Value))
		}
		labels := strings.Builder{}
		for _, v := range nodeGroup.Labels {
			if labels.Len() > 0 {
				labels.WriteString("\n")
			}
			_, _ = labels.WriteString(fmt.Sprintf("%s=%s", v.Key, v.Value))
		}
		taints := strings.Builder{}
		for _, v := range nodeGroup.Taints {
			if taints.Len() > 0 {
				taints.WriteString("\n")
			}
			_, _ = taints.WriteString(fmt.Sprintf("%s=%s:%s", v.Key, v.Value, v.Effect))
		}

		var storageName string
		if storage, err := svc.GetStorageDetails(exec.Context(), &request.GetStorageDetailsRequest{UUID: nodeGroup.Storage}); err != nil {
			storageName = ""
		} else {
			storageName = storage.Title
		}

		nodeGroups = append(nodeGroups, output.CombinedSection{
			Contents: output.Combined{
				output.CombinedSection{
					Contents: output.Details{
						Sections: []output.DetailSection{
							{
								Title: fmt.Sprintf("Node group %d (%s):", i+1, nodeGroup.Name),
								Rows: []output.DetailRow{
									{Title: "Name:", Value: nodeGroup.Name},
									{Title: "Count:", Value: nodeGroup.Count},
									{Title: "Plan:", Value: nodeGroup.Plan},
									{Title: "State:", Value: nodeGroup.State, Format: format.KubernetesNodeGroupState},
									{Title: "Storage UUID:", Value: nodeGroup.Storage, Colour: ui.DefaultUUUIDColours},
									{Title: "Storage name:", Value: storageName, Format: format.PossiblyUnknownString},
									{Title: "Kubelet args:", Value: kubeletArgs.String()},
									{Title: "Labels:", Value: labels.String()},
									{Title: "Taints:", Value: taints.String()},
								},
							},
						},
					},
				},
			},
		})
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: cluster,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: []output.DetailSection{
						{
							Title: "Overview:",
							Rows: []output.DetailRow{
								{Title: "UUID:", Value: cluster.UUID, Colour: ui.DefaultUUUIDColours},
								{Title: "Name:", Value: cluster.Name},
								{Title: "Network UUID:", Value: cluster.Network, Colour: ui.DefaultUUUIDColours},
								{Title: "Network name:", Value: networkName, Format: format.PossiblyUnknownString},
								{Title: "Network CIDR:", Value: cluster.NetworkCIDR, Colour: ui.DefaultAddressColours},
								{Title: "Private node groups:", Value: cluster.PrivateNodeGroups, Format: format.Boolean},
								{Title: "Zone:", Value: cluster.Zone},
								{Title: "Operational state:", Value: cluster.State, Format: format.KubernetesClusterState},
							},
						},
					},
				},
			},
			output.CombinedSection{
				Contents: nodeGroups,
			},
		},
	}, nil
}
