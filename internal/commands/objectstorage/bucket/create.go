package bucket

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CreateCommand creates the 'objectstorage bucket create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a bucket in managed object storage service",
			"upctl object-storage bucket create --service 012345...789 --name my-bucket",
			"upctl object-storage bucket create --service my-service --name my-bucket",
		),
	}
}

type createCommand struct {
	*commands.BaseCommand
	params request.CreateManagedObjectStorageBucketRequest
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Name, "name", "", "The name of the bucket.")
	fs.StringVar(&s.params.ServiceUUID, "service", "", "Service UUID.")
}

// ExecuteWithoutArguments implements Command.ExecuteWithoutArguments
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.params.ServiceUUID == "" {
		return nil, fmt.Errorf("service is required")
	}
	if s.params.Name == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

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
		{Title: "Total Size Bytes", Value: res.TotalSizeBytes},
	}}, nil
}
