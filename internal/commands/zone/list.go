package zone

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
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
	zones, err := svc.GetZones(exec.Context())
	if err != nil {
		return nil, err
	}

	hasParentZones := false
	rows := []output.TableRow{}
	for _, z := range zones.Zones {
		if len(z.ParentZone) > 0 {
			hasParentZones = true
		}

		rows = append(rows, output.TableRow{
			z.ID,
			z.Description,
			z.Public,
			z.ParentZone,
		})
	}

	columns := []output.TableColumn{
		{Key: "id", Header: "ID"},
		{Key: "description", Header: "Description"},
		{Key: "public", Header: "Public", Format: format.Boolean},
	}

	if hasParentZones {
		columns = append(columns, output.TableColumn{Key: "parent_zone", Header: "Parent zone"})
	}

	return output.MarshaledWithHumanOutput{
		Value: zones,
		Output: output.Table{
			Columns: columns,
			Rows:    rows,
		},
	}, nil
}
