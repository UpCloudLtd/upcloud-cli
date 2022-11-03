package ipaddress

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"
)

// ListCommand creates the "ip-address list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List IP addresses", "upctl ip-address list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.MakeExecuteCommand
func (s *listCommand) InitCommand() {
	// TODO: reimplement
	//	flags := &pflag.FlagSet{}
	//	s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
	//	s.AddFlags(flags)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	ips, err := exec.IPAddress().GetIPAddresses()
	if err != nil {
		return nil, err
	}
	rows := make([]output.TableRow, len(ips.IPAddresses))
	for i, ipAddress := range ips.IPAddresses {
		rows[i] = output.TableRow{ipAddress.Address, ipAddress.Access, ipAddress.Family, ipAddress.PartOfPlan, ipAddress.PTRRecord, ipAddress.ServerUUID, ipAddress.Floating, ipAddress.Zone}
	}
	return output.Table{
		Columns: []output.TableColumn{
			{
				Header: "Address",
				Key:    "address",
				Colour: ui.DefaultAddressColours,
			},
			{
				Header: "Access",
				Key:    "access",
			},
			{
				Header: "Family",
				Key:    "family",
			},
			{
				Header: "Part of Plan",
				Key:    "partofplan",
				Format: format.Boolean,
			},
			{
				Header: "PTR Record",
				Key:    "ptrrecord",
			},
			{
				Header: "Server",
				Key:    "server",
				Colour: ui.DefaultUUUIDColours,
			},
			{
				Header: "Floating",
				Key:    "floating",
				Format: format.Boolean,
			},
			{
				Header: "Zone",
				Key:    "zone",
			},
		},
		Rows: rows,
	}, nil
}
