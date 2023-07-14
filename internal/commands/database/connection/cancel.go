package databaseconnection

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/pflag"
)

type cancelCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database
	pid       int
	terminate config.OptionalBoolean
}

// CancelCommand creates the "connection cancel" command
func CancelCommand() commands.Command {
	return &cancelCommand{
		BaseCommand: commands.New(
			"cancel",
			"Terminate client connection or cancel running query for a database",
			`upctl database connection cancel 0fa980c4-0e4f-460b-9869-11b7bd62b833 --pid 2345421`,
			`upctl database connection cancel 0fa980c4-0e4f-460b-9869-11b7bd62b833 --pid 2345421 --terminate`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *cancelCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.IntVar(&s.pid, "pid", 0, "Process ID of the connection to cancel.")
	config.AddToggleFlag(flagSet, &s.terminate, "terminate", false, "Request immediate termination instead of soft cancel.")

	s.AddFlags(flagSet)
	_ = s.Cobra().MarkFlagRequired("pid")
}

// Execute implements commands.MultipleArgumentCommand
func (s *cancelCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Cancelling connection %v to database %v", s.pid, uuid)
	exec.PushProgressStarted(msg)

	if err := svc.CancelManagedDatabaseConnection(exec.Context(), &request.CancelManagedDatabaseConnection{ //nolint:staticcheck // Deprecated, replace in a feature PR
		UUID:      uuid,
		Pid:       s.pid,
		Terminate: s.terminate.Value(),
	}); err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
