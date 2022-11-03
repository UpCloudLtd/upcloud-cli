package ipaddress

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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
	ipAddress, err := exec.IPAddress().GetIPAddressDetails(&request.GetIPAddressDetailsRequest{
		Address: arg,
	})
	if err != nil {
		return nil, err
	}
	return output.Details{
		Sections: []output.DetailSection{
			{
				Rows: []output.DetailRow{
					{Title: "Address:", Key: "address", Value: ipAddress.Address, Colour: ui.DefaultAddressColours},
					{Title: "Access:", Key: "access", Value: ipAddress.Access},
					{Title: "Family:", Key: "access", Value: ipAddress.Family},
					{Title: "Part of Plan:", Key: "access", Value: ipAddress.PartOfPlan, Format: format.Boolean},
					{Title: "PTR Record:", Key: "access", Value: ipAddress.PTRRecord},
					{
						Title: "Server UUID:",
						Key:   "access", Value: ipAddress.ServerUUID,
						Colour: ui.DefaultUUUIDColours,
					},
					{Title: "MAC:", Key: "credits", Value: ipAddress.MAC},
					{Title: "Floating:", Key: "credits", Value: ipAddress.Floating, Format: format.Boolean},
					{Title: "Zone:", Key: "zone", Value: ipAddress.Zone},
				},
			},
		},
	}, nil
}
