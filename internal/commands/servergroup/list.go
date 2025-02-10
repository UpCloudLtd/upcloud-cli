package servergroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the "servergroup list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current server groups", "upctl server-group list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *listCommand) InitCommand() {
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"servergroup"})
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"servergroup"}, "server-group")

	svc := exec.All()
	serverGroups, err := svc.GetServerGroups(exec.Context(), &request.GetServerGroupsRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, serverGroup := range serverGroups {
		groupStatus := notApplicable

		if serverGroup.AntiAffinityPolicy != upcloud.ServerGroupAntiAffinityPolicyOff {
			for _, serverState := range serverGroup.AntiAffinityStatus {
				if groupStatus != unMet {
					groupStatus = string(serverState.Status)
				}
			}
		}

		rows = append(rows, output.TableRow{
			serverGroup.UUID,
			serverGroup.Title,
			serverGroup.AntiAffinityPolicy,
			groupStatus,
			len(serverGroup.Members),
		})
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: serverGroups,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "title", Header: "Title"},
				{Key: "anti_affinity", Header: "Anti-affinity policy"},
				{Key: "anti_affinity_status", Header: "Anti-affinity status", Format: format.ServerGroupAntiAffinityState},
				{Key: "server_count", Header: "Server count"},
			},
			Rows: rows,
		},
	}, nil
}
