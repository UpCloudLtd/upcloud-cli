package user

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CreateAccessKeyCommand creates the 'objectstorage user create-access-key' command
func CreateAccessKeyCommand() commands.Command {
	return &createAccessKeyCommand{
		BaseCommand: commands.New(
			"create-access-key",
			"Create an access key for a user in managed object storage service",
			"upctl object-storage user create-access-key --service 012345...789 --username myuser",
			"upctl object-storage user create-access-key --service my-service --username myuser",
		),
	}
}

type createAccessKeyCommand struct {
	*commands.BaseCommand
	params request.CreateManagedObjectStorageUserAccessKeyRequest
}

// InitCommand implements Command.InitCommand
func (s *createAccessKeyCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`Create an access key for a user in managed object storage service

Creates a new access key for the specified user in a managed object storage service. The access key can be used to authenticate API requests to the object storage service. Note that the secret access key is only shown once during creation and cannot be retrieved later.`)

	fs := s.Cobra().Flags()

	fs.StringVar(&s.params.ServiceUUID, "service", "", "Service UUID or name where the user exists.")
	fs.StringVar(&s.params.Username, "username", "", "The username to create the access key for.")

	commands.Must(s.Cobra().MarkFlagRequired("service"))
	commands.Must(s.Cobra().MarkFlagRequired("username"))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createAccessKeyCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.params.ServiceUUID == "" {
		return nil, fmt.Errorf("service UUID or name is required")
	}

	if s.params.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	svc := exec.All()

	msg := fmt.Sprintf("Creating access key for user %v in service %v", s.params.Username, s.params.ServiceUUID)
	exec.PushProgressStarted(msg)

	res, err := svc.CreateManagedObjectStorageUserAccessKey(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	// Handle SecretAccessKey which might be a pointer
	var secretKey interface{}
	if res.SecretAccessKey != nil {
		secretKey = *res.SecretAccessKey
	} else {
		secretKey = ""
	}

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "Access Key ID", Value: res.AccessKeyID, Colour: ui.DefaultUUUIDColours},
		{Title: "Secret Access Key", Value: secretKey, Colour: ui.DefaultErrorColours},
		{Title: "Status", Value: res.Status},
		{Title: "Created At", Value: res.CreatedAt},
	}}, nil
}
