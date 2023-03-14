package database

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
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
	exec.PushProgressStarted(msg)

	err := exec.All().DeleteManagedDatabase(exec.Context(), &request.DeleteManagedDatabaseRequest{UUID: arg})
	if err != nil {
		return commands.HandleError(exec, fmt.Sprintf("%s: failed", msg), err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
