package storage

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type createBackupCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	completion.Storage
	params createBackupParams
}

type createBackupParams struct {
	request.CreateBackupRequest
}

// CreateBackupCommand creates the "storage backup create" command
func CreateBackupCommand() commands.Command {
	return &createBackupCommand{
		BaseCommand: commands.New("create", "Create backup of a storage"),
	}
}

var defaultCreateBackupParams = &createBackupParams{
	CreateBackupRequest: request.CreateBackupRequest{},
}

// InitCommand implements Command.InitCommand
func (s *createBackupCommand) InitCommand() {
	s.params = createBackupParams{CreateBackupRequest: request.CreateBackupRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", defaultCreateBackupParams.Title, "A short, informational description.")

	s.AddFlags(flagSet)
}

// Execute implements command.Command
func (s *createBackupCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()

	if s.params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	msg := fmt.Sprintf("backing up %v", uuid)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	req := s.params.CreateBackupRequest
	req.UUID = uuid

	res, err := svc.CreateBackup(&req)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
