package servergroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// ListCommand creates the "servergroup list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current server groups", "upctl servergroup list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	clusters, err := svc.GetServerGroups(exec.Context(), &request.GetServerGroupsRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, serverGroup := range clusters {
		rows = append(rows, output.TableRow{
			serverGroup.UUID,
			serverGroup.Title,
			serverGroup.AntiAffinityPolicy,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
			{Key: "title", Header: "Title"},
			{Key: "anti_affinity", Header: "Anti-affinity policy"},
		},
		Rows: rows,
	}, nil
}
