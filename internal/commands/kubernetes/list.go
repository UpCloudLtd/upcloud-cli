package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// ListCommand creates the "kubernetes list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current Kubernetes clusters", "upctl kubernetes list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	clusters, err := svc.GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, cluster := range clusters {
		rows = append(rows, output.TableRow{
			cluster.UUID,
			cluster.Name,
			cluster.Network,
			cluster.NetworkCIDR,
			cluster.Zone,
			cluster.State,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: clusters,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "network", Header: "Network UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "network_cidr", Header: "Network CIDR", Colour: ui.DefaultAddressColours},
				{Key: "zone", Header: "Zone"},
				{Key: "state", Header: "Operational state", Format: format.KubernetesClusterState},
			},
			Rows: rows,
		},
	}, nil
}
