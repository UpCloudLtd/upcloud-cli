package user

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	params createParams
}

type createParams struct {
	request.CreateManagedObjectStorageUserRequest
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.params.Username, "username", "", "Username.")
	fs.StringVar(&s.params.ServiceUUID, "service", "", "Service UUID.")
	s.AddFlags(fs)
}

// CreateCommand creates the 'objectstorage user create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a user in managed object storage service",
			"upctl object-storage user create --service 012345...789 --username myuser",
			"upctl object-storage user create --service my-service --username myuser",
		),
	}
}

func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.params.ServiceUUID == "" {
		return nil, fmt.Errorf("service is required")
	}

	if s.params.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	svc := exec.All()

	serviceUUID := s.params.ServiceUUID

	msg := "Creating user " + s.params.Username + " in service " + serviceUUID
	exec.PushProgressStarted(msg)

	// Build the request
	req := &request.CreateManagedObjectStorageUserRequest{
		ServiceUUID: serviceUUID,
		Username:    s.params.Username,
	}

	exec.PushProgressUpdateMessage(msg, msg)
	res, err := svc.CreateManagedObjectStorageUser(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "Username", Value: res.Username},
		{Title: "ARN", Value: res.ARN},
		{Title: "Created At", Value: res.CreatedAt},
	}}, nil
}
