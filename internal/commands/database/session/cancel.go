package databasesession

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cancelCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database
	pid       int
	terminate config.OptionalBoolean
}

// CancelCommand creates the "session cancel" command
func CancelCommand() commands.Command {
	return &cancelCommand{
		BaseCommand: commands.New(
			"cancel",
			"Terminate client session or cancel running query for a database",
			`upctl database session cancel 0fa980c4-0e4f-460b-9869-11b7bd62b832 --pid 2345422`,
			`upctl database session cancel mysql-1 --pid 2345422 --terminate`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *cancelCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.IntVar(&s.pid, "pid", 0, "Process ID of the session to cancel.")
	config.AddToggleFlag(flagSet, &s.terminate, "terminate", false, "Request immediate termination instead of soft cancel.")

	s.AddFlags(flagSet)
	commands.Must(s.Cobra().MarkFlagRequired("pid"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("pid", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *cancelCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	db, err := svc.GetManagedDatabase(exec.Context(), &request.GetManagedDatabaseRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	switch db.Type {
	case upcloud.ManagedDatabaseServiceTypeMySQL:
		break
	case upcloud.ManagedDatabaseServiceTypePostgreSQL:
		break
	default:
		return nil, fmt.Errorf("session cancel not supported for database type %s", db.Type)
	}

	if db.State != upcloud.ManagedDatabaseStateRunning {
		return nil, fmt.Errorf("database is not in running state")
	}

	msg := fmt.Sprintf("Cancelling session %v to database %v", s.pid, uuid)
	exec.PushProgressStarted(msg)

	if err := svc.CancelManagedDatabaseSession(exec.Context(), &request.CancelManagedDatabaseSession{
		UUID:      uuid,
		Pid:       s.pid,
		Terminate: s.terminate.Value(),
	}); err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
