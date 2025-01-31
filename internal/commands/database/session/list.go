package databasesession

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type listCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database
	limit  int
	offset int
	order  string
}

// ListCommand creates the "session list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List current sessions to specified database",
			"upctl database session list 0fa980c4-0e4f-460b-9869-11b7bd62b832",
			"upctl database session list mysql-1 --limit 16 --offset 32 --order pid:desc",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.IntVar(&s.limit, "limit", 10, "Number of entries to receive at most.")
	flagSet.IntVar(&s.limit, "offset", 0, "Offset for retrieved results based on sort order.")
	flagSet.StringVar(&s.order, "order", "query_duration:desc", "Key and direction for sorting.")

	s.AddFlags(flagSet)
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("limit", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("offset", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("order", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *listCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	db, err := svc.GetManagedDatabase(exec.Context(), &request.GetManagedDatabaseRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	var outputFn func(sessions upcloud.ManagedDatabaseSessions) output.Output

	switch db.Type {
	case upcloud.ManagedDatabaseServiceTypeMySQL:
		outputFn = mysql
	case upcloud.ManagedDatabaseServiceTypePostgreSQL:
		outputFn = pg
	case upcloud.ManagedDatabaseServiceTypeRedis: //nolint:staticcheck // To be removed when Redis support has been removed
		outputFn = redis
	default:
		return nil, fmt.Errorf("session list not supported for database type %s", db.Type)
	}

	if db.State != upcloud.ManagedDatabaseStateRunning {
		return nil, fmt.Errorf("database is not in running state")
	}

	sessions, err := svc.GetManagedDatabaseSessions(exec.Context(), &request.GetManagedDatabaseSessionsRequest{
		Limit:  s.limit,
		Offset: s.offset,
		Order:  s.order,
		UUID:   uuid,
	})
	if err != nil {
		return nil, err
	}

	return outputFn(sessions), nil
}

func mysql(sessions upcloud.ManagedDatabaseSessions) output.Output {
	rows := make([]output.TableRow, 0)

	for _, session := range sessions.MySQL {
		rows = append(rows, output.TableRow{
			session.Id,
			session.Query,
			session.Usename,
			session.ClientAddr,
			session.ApplicationName,
			session.Datname,
			session.QueryDuration.String(),
			session.State,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: sessions,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "id", Header: "Process ID", Format: format.DatabaseSessionPID},
				{Key: "query", Header: "Query"},
				{Key: "usename", Header: "Username"},
				{Key: "client_addr", Header: "Client IP", Colour: ui.DefaultAddressColours},
				{Key: "application_name", Header: "Application"},
				{Key: "datname", Header: "Database"},
				{Key: "query_duration", Header: "Age"},
				{Key: "state", Header: "State", Format: format.DatabaseSessionState},
			},
			Rows: rows,
		},
	}
}

func pg(sessions upcloud.ManagedDatabaseSessions) output.Output {
	rows := make([]output.TableRow, 0)

	for _, session := range sessions.PostgreSQL {
		rows = append(rows, output.TableRow{
			session.Id,
			session.Query,
			session.Usename,
			session.ClientAddr,
			session.ApplicationName,
			session.Datname,
			session.QueryDuration.String(),
			session.State,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: sessions,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "id", Header: "Process ID", Format: format.DatabaseSessionPID},
				{Key: "query", Header: "Query"},
				{Key: "usename", Header: "Username"},
				{Key: "client_addr", Header: "Client IP", Colour: ui.DefaultAddressColours},
				{Key: "application_name", Header: "Application"},
				{Key: "datname", Header: "Database"},
				{Key: "query_duration", Header: "Query duration"},
				{Key: "state", Header: "State", Format: format.DatabaseSessionState},
			},
			Rows: rows,
		},
	}
}

func redis(sessions upcloud.ManagedDatabaseSessions) output.Output {
	rows := make([]output.TableRow, 0)

	for _, session := range sessions.Redis { //nolint:staticcheck // To be removed when Redis support has been removed
		rows = append(rows, output.TableRow{
			session.Id,
			session.Query,
			session.FlagsRaw,
			session.ClientAddr,
			session.ApplicationName,
			session.ActiveDatabase,
			session.ConnectionAge.String(),
			session.ConnectionIdle.String(),
		})
	}
	return output.MarshaledWithHumanOutput{
		Value: sessions,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "id", Header: "Process ID", Format: format.DatabaseSessionPID},
				{Key: "query", Header: "Query"},
				{Key: "flags_raw", Header: "Flags"},
				{Key: "client_addr", Header: "Client IP", Colour: ui.DefaultAddressColours},
				{Key: "application_name", Header: "Application name"},
				{Key: "active_database", Header: "Database"},
				{Key: "connection_age", Header: "Age"},
				{Key: "connection_idle", Header: "Idle"},
			},
			Rows: rows,
		},
	}
}
