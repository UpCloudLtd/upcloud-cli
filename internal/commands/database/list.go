package database

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/paging"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	upcloudv9 "github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
	"github.com/spf13/pflag"
)

// ListCommand creates the "database list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current databases", "upctl database list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	paging.PageParameters
}

func (s *listCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	s.ConfigureFlags(fs)
	s.AddFlags(fs)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	client, err := v9Client(exec)
	if err != nil {
		return nil, err
	}

	limit, offset := s.LimitOffset()
	res, err := client.ListDatabasesWithResponse(exec.Context(), &upcloudv9.ListDatabasesParams{
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, apierror.FromResponse(res.StatusCode(), res.Body)
	}
	databases := *res.JSON200
	legacyOutput, err := legacyDatabaseListOutput(databases)
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, db := range databases {
		title := deref(db.Title)
		if title == "" {
			title = deref(db.Name)
		}

		rows = append(rows, output.TableRow{
			uuidString(db.Uuid),
			title,
			deref(db.Type),
			deref(db.Plan),
			deref(db.Zone),
			deref(db.State),
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: legacyOutput,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "title", Header: "Title"},
				{Key: "type", Header: "Type"},
				{Key: "plan", Header: "Plan"},
				{Key: "zone", Header: "Zone"},
				{Key: "state", Header: "State", Format: format.DatabaseState},
			},
			Rows: rows,
		},
	}, nil
}
