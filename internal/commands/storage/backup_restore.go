package storage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
		BaseCommand: commands.New("restore", "Restore backup of a storage", ""),
	}
}

// InitCommand implements Command.InitCommand
func (s *restoreBackupCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *restoreBackupCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	msg := fmt.Sprintf("restoring backup %v", uuid)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	svc := exec.Storage()
	req := s.params.RestoreBackupRequest
	req.UUID = uuid

	err := svc.RestoreBackup(&req)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.None{}, nil
}
