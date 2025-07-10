package user

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the 'objectstorage user list' command
func ListCommand() commands.Command {
	return &listUsersCommand{
		BaseCommand: commands.New(
			"list",
			"List users in a managed object storage service",
			"upctl object-storage user list <service-uuid>",
			"upctl object-storage user list my-service",
		),
	}
}

type listUsersCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
}

// Execute implements commands.MultipleArgumentCommand
func (s *listUsersCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Listing users in service %s", serviceUUID)
	exec.PushProgressStarted(msg)

	// Build the request
	req := &request.GetManagedObjectStorageUsersRequest{
		ServiceUUID: serviceUUID,
	}

	exec.PushProgressUpdateMessage(msg, msg)
	res, err := svc.GetManagedObjectStorageUsers(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	rows := []output.TableRow{}
	for _, user := range res {
		rows = append(rows, output.TableRow{
			user.Username,
			user.ARN,
			user.CreatedAt.String(),
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: res,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "username", Header: "Username"},
				{Key: "arn", Header: "ARN"},
				{Key: "created_at", Header: "Created"},
			},
			Rows: rows,
		},
	}, nil
}
