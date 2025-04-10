package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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

func (s *showCommand) InitCommand() {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"uks"})
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"uks"}, "k8s")

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

	nodeGroupRows := []output.TableRow{}
	for _, nodeGroup := range cluster.NodeGroups {
		nodeGroupRows = append(nodeGroupRows, output.TableRow{
			nodeGroup.Name,
			nodeGroup.Count,
			nodeGroup.Plan,
			nodeGroup.AntiAffinity,
			nodeGroup.UtilityNetworkAccess,
			nodeGroup.State,
		})
	}

	nodeGroupColumns := []output.TableColumn{
		{Key: "name", Header: "Name"},
		{Key: "count", Header: "Count"},
		{Key: "plan", Header: "Plan"},
		{Key: "anti_affinity", Header: "Anti affinity", Format: format.Boolean},
		{Key: "utility_network_access", Header: "Utility network access", Format: format.Boolean},
		{Key: "state", Header: "State", Format: format.KubernetesNodeGroupState},
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
								{Title: "Version:", Value: cluster.Version},
								{Title: "Network UUID:", Value: cluster.Network, Colour: ui.DefaultUUUIDColours},
								{Title: "Network name:", Value: networkName, Format: format.PossiblyUnknownString},
								{Title: "Network CIDR:", Value: cluster.NetworkCIDR, Colour: ui.DefaultAddressColours},
								{Title: "Kubernetes API allowed IPs:", Value: cluster.ControlPlaneIPFilter, Format: format.IPFilter},
								{Title: "Private node groups:", Value: cluster.PrivateNodeGroups, Format: format.Boolean},
								{Title: "Zone:", Value: cluster.Zone},
								{Title: "Operational state:", Value: cluster.State, Format: format.KubernetesClusterState},
							},
						},
					},
				},
			},
			labels.GetLabelsSectionWithResourceType(cluster.Labels, "cluster"),
			output.CombinedSection{
				Key:   "node_groups",
				Title: "Node groups:",
				Contents: output.Table{
					Columns:      nodeGroupColumns,
					Rows:         nodeGroupRows,
					EmptyMessage: "No node groups found for this cluster.",
				},
			},
		},
	}, nil
}
