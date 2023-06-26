package servergroup

import (
	"fmt"
	"strconv"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

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
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	serverGroup, err := svc.GetServerGroup(exec.Context(), &request.GetServerGroupRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	antiAffinityStatus := output.Combined{}

	for i, memberStatus := range serverGroup.AntiAffinityStatus {
		var serverTitle, serverHostname, serverHost, serverZone string
		if serverDetails, err := svc.GetServerDetails(exec.Context(), &request.GetServerDetailsRequest{UUID: memberStatus.ServerUUID}); err != nil {
			serverTitle = ""
			serverHostname = ""
			serverHost = ""
			serverZone = ""
		} else {
			serverTitle = serverDetails.Title
			serverHostname = serverDetails.Hostname
			serverHost = strconv.Itoa(serverDetails.Host)
			serverZone = serverDetails.Zone
		}

		antiAffinityStatus = append(antiAffinityStatus, output.CombinedSection{
			Contents: output.Combined{
				output.CombinedSection{
					Contents: output.Details{
						Sections: []output.DetailSection{
							{
								Title: fmt.Sprintf("Server %d (%s):", i+1, memberStatus.ServerUUID),
								Rows: []output.DetailRow{
									{Title: "Title:", Value: serverTitle},
									{Title: "Hostname:", Value: serverHostname},
									{Title: "Host:", Value: serverHost},
									{Title: "Zone:", Value: serverZone},
									{Title: "Anti-affinity state:", Value: string(memberStatus.Status), Format: format.ServerGroupAntiAffinityState},
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
							},
						},
					},
				},
			},
			labels.GetLabelsSection(serverGroup.Labels),
			output.CombinedSection{
				Contents: antiAffinityStatus,
			},
		},
	}, nil
}
