package objectstorage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// DeleteCommand creates the "objectstorage delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a Managed object storage service",
			"upctl objectstorage delete 55199a44-4751-4e27-9394-7c7661910be8",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting object storage service %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteManagedObjectStorage(exec.Context(), &request.DeleteManagedObjectStorageRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
