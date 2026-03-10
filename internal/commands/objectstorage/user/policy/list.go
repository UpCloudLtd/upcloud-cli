package policy

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the 'object-storage user-policy list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List policies attached to a user in managed object storage service",
			"upctl object-storage user-policy list <service-uuid> --username myuser",
			"upctl object-storage user-policy list my-service --username myuser",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.GetManagedObjectStorageUserPoliciesRequest
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Username, "username", "", "The username to list policies for.")
	commands.Must(s.Cobra().MarkFlagRequired("username"))
}

// Execute implements Command.Execute
func (s *listCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	s.params.ServiceUUID = serviceUUID
	svc := exec.All()

	msg := fmt.Sprintf("Listing policies for user %s in service %s", s.params.Username, serviceUUID)
	exec.PushProgressStarted(msg)

	res, err := svc.GetManagedObjectStorageUserPolicies(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	rows := []output.TableRow{}
	for _, policy := range res {
		rows = append(rows, output.TableRow{
			policy.Name,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: res,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Policy Name"},
			},
			Rows: rows,
		},
	}, nil
}
