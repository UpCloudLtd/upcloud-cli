package databaseproperties

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/spf13/pflag"
)

// ShowCommand creates the "database properties show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show database property details", "upctl database properties show --db=pg version", "upctl database properties --db=pg version pg_stat_statements_track"),
	}
}

type showCommand struct {
	*commands.BaseCommand
	dbServiceType string
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.dbServiceType, "db", "", "Database service type")
	s.AddFlags(flags)
	_ = s.Cobra().MarkFlagRequired("db")
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, key string) (output.Output, error) {
	svc := exec.All()
	dbType, err := svc.GetManagedDatabaseServiceType(exec.Context(), &request.GetManagedDatabaseServiceTypeRequest{Type: s.dbServiceType})
	if err != nil {
		return nil, err
	}

	details, ok := dbType.Properties[key]
	if !ok {
		return nil, fmt.Errorf(`no property "%s" available for %s database`, key, s.dbServiceType)
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
