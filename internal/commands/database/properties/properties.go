package databaseproperties

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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
	dbType, err := svc.GetManagedDatabaseServiceType(&request.GetManagedDatabaseServiceTypeRequest{Type: serviceType})
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
