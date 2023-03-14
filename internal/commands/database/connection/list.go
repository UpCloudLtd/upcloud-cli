package databaseconnection

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// ListCommand creates the "connection list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current connections to specified databases", "upctl database connection list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database
}

// Execute implements commands.MultipleArgumentCommand
func (s *listCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	connections, err := svc.GetManagedDatabaseConnections(exec.Context(), &request.GetManagedDatabaseConnectionsRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, conn := range connections {
		rows = append(rows, output.TableRow{
			conn.Pid,
			conn.State,
			conn.ApplicationName,
			conn.DatName,
			conn.Username,
			conn.ClientAddr,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: connections,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "pid", Header: "Process ID"},
				{Key: "state", Header: "State"},
				{Key: "application_name", Header: "Application name"},
				{Key: "dataname", Header: "Database"},
				{Key: "username", Header: "Username"},
				{Key: "client_addr", Header: "Client IP", Colour: ui.DefaultAddressColours},
			},
			Rows: rows,
		},
	}, nil
}
