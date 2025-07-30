package accesskey

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the 'object-storage access-key delete' command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete an access key from a user in managed object storage service",
			"upctl object-storage access-key delete <service-uuid> --username myuser --access-key-id AKIAIOSFODNN7EXAMPLE",
			"upctl object-storage access-key delete my-service --username myuser --access-key-id AKIAIOSFODNN7EXAMPLE",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.DeleteManagedObjectStorageUserAccessKeyRequest
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`Delete an access key from a user in managed object storage service\n\nDeletes the specified access key from the user in the managed object storage service. The access key will be permanently deleted and can no longer be used for authentication.`)

	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Username, "username", "", "Username that owns the access key")
	fs.StringVar(&s.params.AccessKeyID, "access-key-id", "", "Access key ID to delete")

	commands.Must(s.Cobra().MarkFlagRequired("username"))
	commands.Must(s.Cobra().MarkFlagRequired("access-key-id"))
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID

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
