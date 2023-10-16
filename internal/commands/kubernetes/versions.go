package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// VersionsCommand creates the "kubernetes versions" command
func VersionsCommand() commands.Command {
	return &versionsCommand{
		BaseCommand: commands.New("versions", "List available versions for Kubernetes clusters", "upctl kubernetes versions"),
	}
}

type versionsCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *versionsCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	versions, err := svc.GetKubernetesVersions(exec.Context(), &request.GetKubernetesVersionsRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, version := range versions {
		rows = append(rows, output.TableRow{
			version,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: versions,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "version", Header: "Version"},
			},
			Rows: rows,
		},
	}, nil
}
