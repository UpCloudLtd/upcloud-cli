package storage

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

// DeleteCommand creates the "storage delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a storage"),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	completion.Storage
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
}

// MaximumExecutions implements command.Command
func (s *deleteCommand) MaximumExecutions() int {
	return maxStorageActions
}

// Execute implements Command.Execute
func (s *deleteCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()
	msg := fmt.Sprintf("deleting storage %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	err := svc.DeleteStorage(&request.DeleteStorageRequest{
		UUID: uuid,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))

	return nil, nil
}
