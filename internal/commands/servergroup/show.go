package servergroup

import (
	"strconv"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// ShowCommand creates the "servergroup show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show server group details",
			"upctl servergroup show 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl servergroup show my-server-group",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingServerGroup
	completion.ServerGroup
}

// Execute implements commands.MultipleArgumentCommand
func (c *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	serverGroup, err := svc.GetServerGroup(exec.Context(), &request.GetServerGroupRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	statusEnabled := false
	groupStatus := notApplicable

	if serverGroup.AntiAffinityPolicy != upcloud.ServerGroupAntiAffinityPolicyOff {
		statusEnabled = true
		groupStatus = met
	}

	statusMap := make(map[string]string, 0)
	if statusEnabled {
		for _, serverState := range serverGroup.AntiAffinityStatus {
			status := string(serverState.Status)
			statusMap[serverState.ServerUUID] = status

			if status == unMet {
				groupStatus = unMet
			}
		}
	}

	serverRows := []output.TableRow{}
	for _, member := range serverGroup.Members {
		serverDetails, err := svc.GetServerDetails(exec.Context(), &request.GetServerDetailsRequest{UUID: member})
		if err != nil {
			return nil, err
		}

		status := groupStatus
		if statusEnabled {
			status = statusMap[member]
		}

		serverRows = append(serverRows, output.TableRow{
			serverDetails.UUID,
			serverDetails.Hostname,
			serverDetails.Zone,
			strconv.Itoa(serverDetails.Host),
			status,
		})
	}

	serverColumns := []output.TableColumn{
		{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
		{Key: "hostname", Header: "Hostname"},
		{Key: "zone", Header: "Zone"},
		{Key: "host", Header: "Host"},
		{Key: "anti_affinity_state", Header: "Anti-affinity state", Format: format.ServerGroupAntiAffinityState},
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
								{Title: "Server count:", Value: len(serverGroup.AntiAffinityStatus)},
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
