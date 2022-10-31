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

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// ShowCommand creates the "kubernetes show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show load balancer details",
			"upctl kubernetes show 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl kubernetes show my-load-balancer",
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
	cluster, err := svc.GetKubernetesCluster(&request.GetKubernetesClusterRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	var networkName string
	if network, err := svc.GetNetworkDetails(&request.GetNetworkDetailsRequest{UUID: cluster.Network}); err != nil {
		networkName = ""
	} else {
		networkName = network.Name
	}

	nodeGroupRows := []output.TableRow{}
	for _, nodeGroup := range cluster.NodeGroups {
		kubeletArgs := make([]string, 0)
		for _, v := range nodeGroup.KubeletArgs {
			kubeletArgs = append(kubeletArgs, fmt.Sprintf("Key: %s\nValue: %s", v.Key, v.Value))
		}
		labels := make([]string, 0)
		for _, v := range nodeGroup.Labels {
			labels = append(labels, fmt.Sprintf("Key: %s\nValue: %s", v.Key, v.Value))
		}
		taints := make([]string, 0)
		for _, v := range nodeGroup.Taints {
			taints = append(taints, fmt.Sprintf("Key: %s\nValue: %s\nEffect: %s", v.Key, v.Value, v.Effect))
		}
		nodeGroupRows = append(nodeGroupRows, output.TableRow{
			nodeGroup.Name,
			nodeGroup.Count,
			nodeGroup.Plan,
			nodeGroup.Storage,
			strings.Join(kubeletArgs, "\n\n"),
			strings.Join(labels, "\n\n"),
			strings.Join(taints, "\n\n"),
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
								{Title: "Network name", Value: networkName, Format: format.PossiblyUnknownString},
								{Title: "Network CIDR:", Value: cluster.NetworkCIDR, Colour: ui.DefaultAddressColours},
								{Title: "Zone", Value: cluster.Zone},
								{Title: "Operational state:", Value: cluster.State, Format: format.KubernetesState},
							},
						},
					},
				},
			},
			output.CombinedSection{
				Title: "Node groups:",
				Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "name", Header: "Name"},
						{Key: "count", Header: "Count"},
						{Key: "plan", Header: "Plan"},
						{Key: "storage", Header: "Storage"},
						{Key: "kubelet_args", Header: "Kubelet args"},
						{Key: "labels", Header: "Labels"},
						{Key: "taints", Header: "Taints"},
					},
					Rows: nodeGroupRows,
				},
			},
		},
	}, nil
}
