package ipaddress

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
)

// ListCommand creates the "ip-address list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List IP addresses", ""),
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
				Color:  ui.DefaultAddressColours,
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
				Format: output.BoolFormat,
			},
			{
				Header: "PTR Record",
				Key:    "ptrrecord",
			},
			{
				Header: "Server",
				Key:    "server",
				Color:  ui.DefaultUUUIDColours,
			},
			{
				Header: "Floating",
				Key:    "floating",
				Format: output.BoolFormat,
			},
			{
				Header: "Zone",
				Key:    "zone",
			},
		},
		Rows: rows,
	}, nil
}
