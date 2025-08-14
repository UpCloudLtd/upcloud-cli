package storagebackup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
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
		BaseCommand: commands.New(
			"create",
			"Create backup of a storage",
			`upctl storage backup create 01cbea5e-eb5b-4072-b2ac-9b635120e5d8 --title "first backup"`,
			`upctl storage backup create "My Storage" --title second_backup`,
		),
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
	commands.Must(s.Cobra().MarkFlagRequired("title"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("title", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *createBackupCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()

	msg := fmt.Sprintf("Backing up storage %v to %v", uuid, s.params.Title)
	exec.PushProgressStarted(msg)

	req := s.params.CreateBackupRequest
	req.UUID = uuid

	res, err := svc.CreateBackup(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
