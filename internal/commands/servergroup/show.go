package servergroup

import (
	"strconv"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ShowCommand creates the "servergroup show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show server group details",
			"upctl server-group show 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl server-group show my-server-group",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingServerGroup
	completion.ServerGroup
}

// InitCommand implements Command.InitCommand
func (c *showCommand) InitCommand() {
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"servergroup"})
}

// Execute implements commands.MultipleArgumentCommand
func (c *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"servergroup"}, "server-group")

	svc := exec.All()
	serverGroup, err := svc.GetServerGroup(exec.Context(), &request.GetServerGroupRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	groupStatus := notApplicable

	statusMap := make(map[string]string, 0)
	if serverGroup.AntiAffinityPolicy != upcloud.ServerGroupAntiAffinityPolicyOff {
		for _, serverState := range serverGroup.AntiAffinityStatus {
			status := string(serverState.Status)
			statusMap[serverState.ServerUUID] = status

			if groupStatus != unMet {
				groupStatus = status
			}
		}
	}

	serverRows := []output.TableRow{}
	for _, serverUUID := range serverGroup.Members {
		serverDetails, err := svc.GetServerDetails(exec.Context(), &request.GetServerDetailsRequest{UUID: serverUUID})
		if err != nil {
			return nil, err
		}

		status, ok := statusMap[serverUUID]
		if !ok {
			status = notApplicable
		}

		host := strconv.Itoa(serverDetails.Host)
		if host == "0" {
			host = notApplicable
		}

		serverRows = append(serverRows, output.TableRow{
			serverDetails.UUID,
			serverDetails.Hostname,
			serverDetails.Zone,
			host,
			status,
			serverDetails.State,
		})
	}

	serverColumns := []output.TableColumn{
		{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
		{Key: "hostname", Header: "Hostname"},
		{Key: "zone", Header: "Zone"},
		{Key: "host", Header: "Host"},
		{Key: "anti_affinity_state", Header: "Anti-affinity state", Format: format.ServerGroupAntiAffinityState},
		{Key: "state", Header: "State", Format: format.ServerState},
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: serverGroup,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: []output.DetailSection{
						{
							Title: "Overview:",
							Rows: []output.DetailRow{
								{Title: "UUID:", Value: serverGroup.UUID, Colour: ui.DefaultUUUIDColours},
								{Title: "Title:", Value: serverGroup.Title},
								{Title: "Anti-affinity policy:", Value: serverGroup.AntiAffinityPolicy},
								{Title: "Anti-affinity state:", Value: groupStatus, Format: format.ServerGroupAntiAffinityState},
								{Title: "Server count:", Value: len(serverGroup.Members)},
							},
						},
					},
				},
			},
			labels.GetLabelsSection(serverGroup.Labels),
			output.CombinedSection{
				Key:   "servers",
				Title: "Servers:",
				Contents: output.Table{
					Columns:      serverColumns,
					Rows:         serverRows,
					EmptyMessage: "No servers in this server group.",
				},
			},
		},
	}, nil
}
