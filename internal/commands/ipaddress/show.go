package ipaddress

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ShowCommand creates the 'ip-address show' command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show current IP address",
			"upctl ip-address show 185.70.196.47",
			"upctl ip-address show 2a04:3544:8000:1000:d40e:4aff:fe6f:5d34",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	completion.IPAddress
	resolver.CachingIPAddress
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
}

func (s *showCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	ipAddress, err := exec.IPAddress().GetIPAddressDetails(exec.Context(), &request.GetIPAddressDetailsRequest{
		Address: arg,
	})
	if err != nil {
		return nil, err
	}

	details := output.Details{
		Sections: []output.DetailSection{
			{
				Rows: []output.DetailRow{
					{Title: "Address:", Key: "address", Value: ipAddress.Address, Colour: ui.DefaultAddressColours},
					{Title: "Access:", Key: "access", Value: ipAddress.Access},
					{Title: "Family:", Key: "family", Value: ipAddress.Family},
					{Title: "Part of Plan:", Key: "part_of_plan", Value: ipAddress.PartOfPlan, Format: format.Boolean},
					{Title: "PTR Record:", Key: "ptr_record", Value: ipAddress.PTRRecord},
					{
						Title: "Server UUID:",
						Key:   "server", Value: ipAddress.ServerUUID,
						Colour: ui.DefaultUUUIDColours,
					},
					{Title: "MAC:", Key: "mac", Value: ipAddress.MAC},
					{Title: "Floating:", Key: "floating", Value: ipAddress.Floating, Format: format.Boolean},
					{Title: "Zone:", Key: "zone", Value: ipAddress.Zone},
				},
			},
		},
	}

	return output.MarshaledWithHumanOutput{
		Value:  ipAddress,
		Output: details,
	}, nil
}
