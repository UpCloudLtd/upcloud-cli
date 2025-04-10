package kubernetes

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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

func (s *versionsCommand) InitCommand() {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"uks"})
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *versionsCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"uks"}, "k8s")

	svc := exec.All()
	versions, err := svc.GetKubernetesVersions(exec.Context(), &request.GetKubernetesVersionsRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, version := range versions {
		rows = append(rows, output.TableRow{
			version.Id,
			version.Version,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: versions,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "id", Header: "ID"},
				{Key: "version", Header: "Version"},
			},
			Rows: rows,
		},
	}, nil
}
