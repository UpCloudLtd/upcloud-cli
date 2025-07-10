package bucket

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the 'objectstorage bucket delete' command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a bucket from a managed object storage service",
			"upctl object-storage bucket delete <service-uuid> --name my-bucket",
			"upctl object-storage bucket delete my-service --name my-bucket",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.DeleteManagedObjectStorageBucketRequest
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Name, "name", "", "The name of the bucket to delete.")
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

// Execute implements Command.Execute
func (s *deleteCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID

	svc := exec.All()

	msg := fmt.Sprintf("Deleting bucket %s from service %s", s.params.Name, serviceUUID)
	exec.PushProgressStarted(msg)

	err := svc.DeleteManagedObjectStorageBucket(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: fmt.Sprintf("Bucket %s deleted from service %s", s.params.Name, serviceUUID)}, nil
}
