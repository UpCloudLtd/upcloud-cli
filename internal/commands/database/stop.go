package database

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.ShutdownManagedDatabase(&request.ShutdownManagedDatabaseRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
