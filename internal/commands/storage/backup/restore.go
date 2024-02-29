package storagebackup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type restoreBackupCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	completion.Storage
	params restoreBackupParams
}

type restoreBackupParams struct {
	request.RestoreBackupRequest
}

// RestoreBackupCommand creates the "storage backup restore" command
func RestoreBackupCommand() commands.Command {
	return &restoreBackupCommand{
		BaseCommand: commands.New(
			"restore",
			"Restore backup of a storage",
			"upctl storage backup restore 01177c9e-7f76-4ce4-b128-bcaa3448f7ec",
			`upctl storage backup restore second_backup`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *restoreBackupCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *restoreBackupCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	msg := fmt.Sprintf("Restoring backup %v", uuid)
	exec.PushProgressSuccess(msg)

	svc := exec.Storage()
	req := s.params.RestoreBackupRequest
	req.UUID = uuid

	err := svc.RestoreBackup(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
