package storage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "storage delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a storage",
			"upctl storage delete 01ac5319-08ac-4e7b-81e5-3140d2bbd7d8",
			"upctl storage delete 0175bb34-8aed-47ce-9290-10cc45f78601 01fcb78f-e73d-4e4d-af5a-0bd6cdba4306",
			`upctl storage delete "My Storage" --backups keep_latest`,
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	completion.Storage

	backups string
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	fs := &pflag.FlagSet{}

	backupsOptions := []string{
		string(request.DeleteStorageBackupsModeKeep),
		string(request.DeleteStorageBackupsModeKeepLatest),
		string(request.DeleteStorageBackupsModeDelete),
	}
	fs.StringVar(&s.backups, "backups", "", "Controls what to do with backups related to the storage. Valid values are "+namedargs.ValidValuesHelp(backupsOptions...)+".")

	s.AddFlags(fs)

	commands.Must(s.Cobra().RegisterFlagCompletionFunc("backups", cobra.FixedCompletions(backupsOptions, cobra.ShellCompDirectiveNoFileComp)))
}

// MaximumExecutions implements command.Command
func (s *deleteCommand) MaximumExecutions() int {
	return maxStorageActions
}

func Delete(exec commands.Executor, uuid string, backups string) (output.Output, error) {
	svc := exec.Storage()
	msg := fmt.Sprintf("Deleting storage %v", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteStorage(exec.Context(), &request.DeleteStorageRequest{
		UUID:    uuid,
		Backups: request.DeleteStorageBackupsMode(backups),
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	return Delete(exec, uuid, s.backups)
}
