package policy

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DetachCommand creates the 'object-storage user-policy detach' command
func DetachCommand() commands.Command {
	return &detachCommand{
		BaseCommand: commands.New(
			"detach",
			"Detach a policy from a user in managed object storage service",
			"upctl object-storage user-policy detach <service-uuid> --username myuser --policy ECSS3FullAccess",
			"upctl object-storage user-policy detach my-service --username myuser --policy ECSS3ReadOnlyAccess",
		),
	}
}

type detachCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.DetachManagedObjectStorageUserPolicyRequest
}

// InitCommand implements Command.InitCommand
func (s *detachCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Username, "username", "", "The username to detach the policy from.")
	fs.StringVar(&s.params.Name, "policy", "", "The name of the policy to detach.")
	commands.Must(s.Cobra().MarkFlagRequired("username"))
	commands.Must(s.Cobra().MarkFlagRequired("policy"))
}

// Execute implements Command.Execute
func (s *detachCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID
	svc := exec.All()

	msg := fmt.Sprintf("Detaching policy %s from user %s in service %s", s.params.Name, s.params.Username, serviceUUID)
	exec.PushProgressStarted(msg)

	err := svc.DetachManagedObjectStorageUserPolicy(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: fmt.Sprintf("Policy %s detached from user %s in service %s", s.params.Name, s.params.Username, serviceUUID)}, nil
}
