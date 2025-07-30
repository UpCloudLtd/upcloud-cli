package user

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the 'objectstorage user delete' command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a user from a managed object storage service",
			"upctl object-storage user delete <service-uuid> --username myuser",
			"upctl object-storage user delete my-service --username myuser",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.DeleteManagedObjectStorageUserRequest
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Username, "username", "", "The username of the user to delete.")
	commands.Must(s.Cobra().MarkFlagRequired("username"))
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *deleteCommand) MaximumExecutions() int {
	return 1
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID

	svc := exec.All()

	msg := fmt.Sprintf("Deleting user %s from service %s", s.params.Username, serviceUUID)
	exec.PushProgressStarted(msg)

	exec.PushProgressUpdateMessage(msg, msg)
	err := svc.DeleteManagedObjectStorageUser(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: fmt.Sprintf("User %s deleted from service %s", s.params.Username, serviceUUID)}, nil
}
