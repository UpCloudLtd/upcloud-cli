package user

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params createParams
}

type createParams struct {
	request.CreateManagedObjectStorageUserRequest
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.params.Username, "username", "", "Username.")
	s.AddFlags(fs)
	commands.Must(s.Cobra().MarkFlagRequired("username"))
}

// CreateCommand creates the 'objectstorage user create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a user in managed object storage service",
			"upctl object-storage user create <service-uuid> --username myuser",
			"upctl object-storage user create my-service --username myuser",
		),
	}
}

// Execute implements Command.Execute
func (s *createCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID

	svc := exec.All()

	msg := fmt.Sprintf("Creating user %s in service %s", s.params.Username, serviceUUID)
	exec.PushProgressStarted(msg)

	res, err := svc.CreateManagedObjectStorageUser(exec.Context(), &s.params.CreateManagedObjectStorageUserRequest)
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
