package database

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// ListCommand creates the "database list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current databases", "upctl database list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	databases, err := svc.GetManagedDatabases(&request.GetManagedDatabasesRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, db := range databases {
		rows = append(rows, output.TableRow{
			db.UUID,
			db.Title,
			db.Type,
			db.Plan,
			db.Zone,
			db.State,
		})
	}

	return output.Table{
		Columns: []output.TableColumn{
			{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
			{Key: "title", Header: "Title"},
			{Key: "type", Header: "Type"},
			{Key: "plan", Header: "Plan"},
			{Key: "zone", Header: "Zone"},
			{Key: "state", Header: "State"},
		},
		Rows: rows,
	}, nil
}
