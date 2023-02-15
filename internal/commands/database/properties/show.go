package databaseproperties

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

// ShowCommand creates the "database properties <serviceType> show" command
func ShowCommand(serviceType string, serviceName string) commands.Command {
	return &showCommand{
		BaseCommand:      commands.New("show", fmt.Sprintf("Show %s database property details", serviceName), fmt.Sprintf("upctl database properties %s show version", serviceType)),
		serviceType:      serviceType,
		DatabaseProperty: completion.DatabaseProperty{ServiceType: serviceType},
	}
}

type showCommand struct {
	*commands.BaseCommand
	completion.DatabaseProperty
	serviceType string
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, key string) (output.Output, error) {
	svc := exec.All()
	dbType, err := svc.GetManagedDatabaseServiceType(exec.Context(), &request.GetManagedDatabaseServiceTypeRequest{Type: s.serviceType})
	if err != nil {
		return nil, err
	}

	details, ok := dbType.Properties[key]
	if !ok {
		return nil, fmt.Errorf(`no property "%s" available for %s database`, key, s.serviceType)
	}

	rows := []output.DetailRow{
		{Title: "Key:", Key: "key", Value: key},
		{Title: "Title:", Key: "title", Value: details.Title},
		{Title: "Description:", Key: "description", Value: details.Description},
		{Title: "Help message:", Key: "user_error", Value: details.UserError},
		{Title: "Create only:", Key: "createOnly", Value: details.CreateOnly, Format: format.Boolean},
		{Title: "Type:", Key: "type", Value: details.Type, Format: formatAlternatives},
		{Title: "Default:", Key: "default", Value: details.Default},
		{Title: "Possible values:", Key: "enum", Value: details.Enum, Format: formatAlternatives},
		{Title: "Pattern:", Key: "pattern", Value: details.Pattern},
		{Title: "Max length:", Key: "maxLength", Value: details.MaxLength},
		{Title: "Min length:", Key: "minLength", Value: details.MinLength},
	}

	return output.Details{
		Sections: []output.DetailSection{
			{
				Rows: filterOutEmptyRows(rows),
			},
		},
	}, nil
}

func filterOutEmptyRows(rows []output.DetailRow) []output.DetailRow {
	nonEmpty := []output.DetailRow{}
	for _, row := range rows {
		if row.Value == nil {
			continue
		}

		if val, ok := row.Value.(string); ok && val == "" {
			continue
		}

		if val, ok := row.Value.(int); ok && val == 0 {
			continue
		}

		nonEmpty = append(nonEmpty, row)
	}

	return nonEmpty
}
