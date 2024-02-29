package database

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// TypesCommand creates the "database types" command
func TypesCommand() commands.Command {
	return &typesCommand{
		BaseCommand: commands.New("types", "List available database types", "upctl database types"),
	}
}

type typesCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *typesCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	dbTypes, err := svc.GetManagedDatabaseServiceTypes(exec.Context(), &request.GetManagedDatabaseServiceTypesRequest{})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, dbType := range dbTypes {
		rows = append(rows, output.TableRow{
			dbType.Name,
			dbType.Description,
			dbType.LatestAvailableVersion,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: dbTypes,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "description", Header: "Description"},
				{Key: "latest_available_version", Header: "Latest Available Version"},
			},
			Rows: rows,
		},
	}, nil
}
