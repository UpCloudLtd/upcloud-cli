package databaseproperties

import (
	"sort"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

// PropertiesCommand creates the "database properties" command
func PropertiesCommand() commands.Command {
	return &propertiesCommand{
		BaseCommand: commands.New("properties", "List available properties for given database type", "upctl database properties pg", "upctl database properties mysql"),
	}
}

type propertiesCommand struct {
	*commands.BaseCommand
}

// Execute implements commands.MultipleArgumentCommand
func (s *propertiesCommand) Execute(exec commands.Executor, serviceType string) (output.Output, error) {
	svc := exec.All()
	dbType, err := svc.GetManagedDatabaseServiceType(exec.Context(), &request.GetManagedDatabaseServiceTypeRequest{Type: serviceType})
	if err != nil {
		return nil, err
	}

	properties := dbType.Properties
	rows := []output.TableRow{}
	for key, details := range properties {
		enumOrExample := details.Enum
		if enumOrExample == nil {
			enumOrExample = details.Example
		}

		rows = append(rows, output.TableRow{
			key,
			details.CreateOnly,
			details.Type,
			enumOrExample,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		iKey := rows[i][0].(string)
		jKey := rows[j][0].(string)
		return iKey < jKey
	})

	return output.MarshaledWithHumanOutput{
		Value: properties,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "property", Header: "Property"},
				{Key: "createOnly", Header: "Create only", Format: format.Boolean},
				{Key: "type", Header: "Type", Format: formatAlternatives},
				{Key: "example", Header: "Example", Format: formatAlternatives},
			},
			Rows: rows,
		},
	}, nil
}
