package databaseproperties

import (
	"fmt"
	"sort"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/utils"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// creates the "database properties <serviceType>" command
func DBTypeCommand(serviceType string, serviceName string) commands.Command {
	return &dbTypeCommand{
		BaseCommand: commands.New(serviceType, fmt.Sprintf("List available properties for %s databases", serviceName), fmt.Sprintf("upctl database properties %s", serviceType)),
		serviceType: serviceType,
	}
}

type dbTypeCommand struct {
	*commands.BaseCommand
	serviceType string
}

// Execute implements commands.NoArgumentCommand
func (s *dbTypeCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	dbType, err := svc.GetManagedDatabaseServiceType(exec.Context(), &request.GetManagedDatabaseServiceTypeRequest{Type: s.serviceType})
	if err != nil {
		return nil, err
	}

	properties := dbType.Properties
	rows := []output.TableRow{}
	for key, details := range utils.GetFlatDatabaseProperties(properties) {
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
				{Key: "type", Header: "Type", Format: format.StringSliceOr},
				{Key: "example", Header: "Example", Format: format.StringSliceOr},
			},
			Rows: rows,
		},
	}, nil
}
