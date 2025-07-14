package accesskey

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CreateCommand creates the 'object-storage access-key create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create an access key for a user in managed object storage service",
			"upctl object-storage access-key create <service-uuid> --username myuser",
			"upctl object-storage access-key create my-service --username myuser",
		),
	}
}

type createCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.CreateManagedObjectStorageUserAccessKeyRequest
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`Create an access key for a user in managed object storage service\n\nCreates a new access key for the specified user in a managed object storage service. The access key can be used to authenticate API requests to the object storage service. Note that the secret access key is only shown once during creation and cannot be retrieved later.`)

	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Username, "username", "", "The username to create the access key for.")
	commands.Must(s.Cobra().MarkFlagRequired("username"))
}

// Execute implements Command.Execute
func (s *createCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID

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
		{Title: "Access Key ID", Value: res.AccessKeyID},
		{Title: "Secret Access Key", Value: secretKey},
		{Title: "Status", Value: res.Status},
		{Title: "Created At", Value: res.CreatedAt},
	}}, nil
}
