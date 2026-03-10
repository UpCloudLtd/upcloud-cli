package filestorage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "file-storage delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a file storage service",
			"upctl file-storage delete 55199a44-4751-4e27-9394-7c7661910be8",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingFileStorage
	completion.FileStorage

	wait config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &c.wait, "wait", false, "Wait until the file storage instance has been deleted before returning.")
	c.AddFlags(flags)
}

func Delete(exec commands.Executor, uuid string, wait bool) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting file storage service %v", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteFileStorage(exec.Context(), &request.DeleteFileStorageRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for file storage service %s to be deleted", uuid))
		err = svc.WaitForFileStorageDeletion(exec.Context(), &request.WaitForFileStorageDeletionRequest{UUID: uuid})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		exec.PushProgressUpdateMessage(msg, msg)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	return Delete(exec, arg, c.wait.Value())
}
