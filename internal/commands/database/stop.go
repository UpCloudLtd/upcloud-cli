package database

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// StopCommand creates the "database stop" command
func StopCommand() commands.Command {
	return &stopCommand{
		BaseCommand: commands.New(
			"stop",
			"Stop a managed database",
			"upctl database stop b0952286-1193-4a81-a1af-62efc014ae4b",
			"upctl database stop b0952286-1193-4a81-a1af-62efc014ae4b 666bcd3c-5c63-428d-a4fd-07c27469a5a6",
			"upctl database stop pg-1x1xcpu-2gb-25gb-pl-waw1",
		),
	}
}

type stopCommand struct {
	*commands.BaseCommand
	completion.Database
	resolver.CachingDatabase
}

// InitCommand implements Command.InitCommand
func (s *stopCommand) InitCommand() {
	s.Cobra().Aliases = []string{"shutdown"}
}

// Execute implements commands.MultipleArgumentCommand
func (s *stopCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Stopping database %v", uuid)
	exec.PushProgressStarted(msg)

	res, err := svc.ShutdownManagedDatabase(exec.Context(), &request.ShutdownManagedDatabaseRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
