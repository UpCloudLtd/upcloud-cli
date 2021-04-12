package ipaddress

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

// ShowCommand creates the 'ip-address show' command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show current IP address"),
	}
}

type showCommand struct {
	*commands.BaseCommand
	completion.IPAddress
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	// TODO: reimplmement
	// s.SetPositionalArgHelp(positionalArgHelp)
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
					{Title: "Address:", Key: "address", Value: ipAddress.Address, Color: ui.DefaultAddressColours},
					{Title: "Access:", Key: "access", Value: ipAddress.Access},
					{Title: "Family:", Key: "access", Value: ipAddress.Family},
					{Title: "Part of Plan:", Key: "access", Value: ipAddress.PartOfPlan, Format: output.BoolFormat},
					{Title: "PTR Record:", Key: "access", Value: ipAddress.PTRRecord},
					{Title: "Server UUID:", Key: "access", Value: ipAddress.ServerUUID},
					{Title: "MAC:", Key: "credits", Value: ipAddress.MAC},
					{Title: "Floating:", Key: "credits", Value: ipAddress.Floating, Format: output.BoolFormat},
					{Title: "Zone:", Key: "zone", Value: ipAddress.Zone},
				},
			},
		},
	}, nil
}
