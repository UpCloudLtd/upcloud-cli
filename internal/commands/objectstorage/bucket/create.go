package bucket

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CreateCommand creates the 'objectstorage bucket create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a bucket in managed object storage service",
			"upctl object-storage bucket create 012345...789 --name my-bucket",
			"upctl object-storage bucket create my-service --name my-bucket",
		),
	}
}

type createCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.CreateManagedObjectStorageBucketRequest
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Name, "name", "", "The name of the bucket.")
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

// ExecuteSingleArgument implements Command.SingleArgumentCommand
func (s *createCommand) ExecuteSingleArgument(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID
	svc := exec.All()

	msg := fmt.Sprintf("Creating bucket %v in service %v", s.params.Name, s.params.ServiceUUID)
	exec.PushProgressStarted(msg)

	res, err := svc.CreateManagedObjectStorageBucket(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "Name", Value: res.Name},
	}}, nil
}
