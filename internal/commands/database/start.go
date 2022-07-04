package database

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// StartCommand creates the "database start" command
func StartCommand() commands.Command {
	return &startCommand{
		BaseCommand: commands.New(
			"start",
			"Start on a managed database",
			"upctl database start b0952286-1193-4a81-a1af-62efc014ae4b",
			"upctl database start b0952286-1193-4a81-a1af-62efc014ae4b 666bcd3c-5c63-428d-a4fd-07c27469a5a6",
			"upctl database start pg-1x1xcpu-2gb-25gb-pl-waw1",
		),
	}
}

type startCommand struct {
	*commands.BaseCommand
	completion.Database
	resolver.CachingDatabase
}

// Execute implements commands.MultipleArgumentCommand
func (s *startCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Starting database %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.StartManagedDatabase(&request.StartManagedDatabaseRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
