package zone

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
)

// ListCommand creates the "zone list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List available zones", "upctl zone list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	zones, err := svc.GetZones()
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, z := range zones.Zones {
		rows = append(rows, output.TableRow{
			z.ID,
			z.Description,
			z.Public,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "id", Header: "ID"},
			{Key: "description", Header: "Description"},
			{Key: "public", Header: "Public", Format: format.Boolean},
		},
		Rows: rows,
	}, nil
}
