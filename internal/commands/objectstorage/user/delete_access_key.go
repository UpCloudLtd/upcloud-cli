package user

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteAccessKeyCommand creates the 'objectstorage user delete-access-key' command
func DeleteAccessKeyCommand() commands.Command {
	return &deleteAccessKeyCommand{
		BaseCommand: commands.New(
			"delete-access-key",
			"Delete an access key from a user in managed object storage service",
			"upctl object-storage delete-access-key <service-uuid> --username myuser --access-key-id AKIAIOSFODNN7EXAMPLE",
			"upctl object-storage delete-access-key my-service --username myuser --access-key-id AKIAIOSFODNN7EXAMPLE",
		),
	}
}

type deleteAccessKeyCommand struct {
	*commands.BaseCommand
	params request.DeleteManagedObjectStorageUserAccessKeyRequest
}

// InitCommand implements Command.InitCommand
func (s *deleteAccessKeyCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`Delete an access key from a user in managed object storage service

Deletes the specified access key from the user in the managed object storage service. The access key will be permanently deleted and can no longer be used for authentication.`)

	fs := s.Cobra().Flags()

	fs.StringVar(&s.params.Username, "username", "", "Username that owns the access key")
	fs.StringVar(&s.params.AccessKeyID, "access-key-id", "", "Access key ID to delete")

	commands.Must(s.Cobra().MarkFlagRequired("username"))
	commands.Must(s.Cobra().MarkFlagRequired("access-key-id"))
}

// MaximumExecutions implements commands.MultipleArgumentCommand
func (s *deleteAccessKeyCommand) MaximumExecutions() int {
	return 1
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteAccessKeyCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID
	if s.params.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	if s.params.AccessKeyID == "" {
		return nil, fmt.Errorf("access key ID is required")
	}

	svc := exec.All()

	msg := fmt.Sprintf("Deleting access key %s for user %s from service %s", s.params.AccessKeyID, s.params.Username, serviceUUID)
	exec.PushProgressStarted(msg)

	err := svc.DeleteManagedObjectStorageUserAccessKey(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: fmt.Sprintf("Access key %s deleted for user %s from service %s", s.params.AccessKeyID, s.params.Username, serviceUUID)}, nil
}
