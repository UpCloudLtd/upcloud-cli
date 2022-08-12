package database

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database
}

// DeleteCommand creates the "delete database" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a database",
			"upctl database delete 0497728e-76ef-41d0-997f-fa9449eb71bc",
			"upctl database delete my_database",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Deleting database %s", arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	err := exec.All().DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: arg})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.None{}, nil
}
