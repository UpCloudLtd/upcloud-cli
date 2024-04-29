package host

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ListCommand creates the "host list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List private cloud hosts", "upctl host list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	hosts, err := svc.GetHosts(exec.Context())
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, h := range hosts.Hosts {
		rows = append(rows, output.TableRow{
			h.ID,
			h.Description,
			h.Zone,
			h.WindowsEnabled,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: hosts,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "id", Header: "ID", Format: formatID},
				{Key: "description", Header: "Description"},
				{Key: "zone", Header: "Zone"},
				{Key: "windows_enabled", Header: "Windows enabled", Format: format.Boolean},
			},
			Rows: rows,
		},
	}, nil
}

func formatID(val interface{}) (text.Colors, string, error) {
	id, ok := val.(int)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse IP addresses from %T, expected int", val)
	}

	return ui.DefaultUUUIDColours, fmt.Sprintf("%d", id), nil
}
