package policy

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// AttachCommand creates the 'object-storage user-policy attach' command
func AttachCommand() commands.Command {
	return &attachCommand{
		BaseCommand: commands.New(
			"attach",
			"Attach a policy to a user in managed object storage service",
			"upctl object-storage user-policy attach <service-uuid> --username myuser --policy ECSS3FullAccess",
			"upctl object-storage user-policy attach my-service --username myuser --policy ECSS3ReadOnlyAccess",
		),
	}
}

type attachCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.AttachManagedObjectStorageUserPolicyRequest
}

// InitCommand implements Command.InitCommand
func (s *attachCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Username, "username", "", "The username to attach the policy to.")
	fs.StringVar(&s.params.Name, "policy", "", "The name of the policy to attach.")
	commands.Must(s.Cobra().MarkFlagRequired("username"))
	commands.Must(s.Cobra().MarkFlagRequired("policy"))
}

// Execute implements Command.Execute
func (s *attachCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID
	svc := exec.All()

	msg := fmt.Sprintf("Attaching policy %s to user %s in service %s", s.params.Name, s.params.Username, serviceUUID)
	exec.PushProgressStarted(msg)

	err := svc.AttachManagedObjectStorageUserPolicy(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: fmt.Sprintf("Policy %s attached to user %s in service %s", s.params.Name, s.params.Username, serviceUUID)}, nil
}
