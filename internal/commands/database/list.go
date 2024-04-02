package database

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/paging"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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
	svc := exec.All()
	databases, err := svc.GetManagedDatabases(exec.Context(), &request.GetManagedDatabasesRequest{
		Page: s.Page(),
	})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, db := range databases {
		title := db.Title
		if title == "" {
			title = db.Name
		}

		rows = append(rows, output.TableRow{
			db.UUID,
			title,
			db.Type,
			db.Plan,
			db.Zone,
			db.State,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: databases,
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
