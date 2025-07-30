package accesskey

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the 'object-storage access-key list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List access keys for a user in managed object storage service",
			"upctl object-storage access-key list <service-uuid> --username myuser",
			"upctl object-storage access-key list my-service --username myuser",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.GetManagedObjectStorageUserAccessKeysRequest
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`List access keys for a user in managed object storage service\n\nLists all access keys for the specified user in the managed object storage service. This helps you find the access key IDs needed for deletion.`)

	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.Username, "username", "", "The username of the user to list access keys from.")
	commands.Must(s.Cobra().MarkFlagRequired("username"))
}

// Execute implements commands.MultipleArgumentCommand
func (s *listCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	s.params.ServiceUUID = uuid

	svc := exec.All()

	msg := fmt.Sprintf("Listing access keys for user %s in service %s", s.params.Username, uuid)
	exec.PushProgressStarted(msg)

	res, err := svc.GetManagedObjectStorageUserAccessKeys(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	rows := []output.TableRow{}
	for _, accessKey := range res {
		rows = append(rows, output.TableRow{
			accessKey.AccessKeyID,
			accessKey.Status,
			accessKey.CreatedAt.String(),
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: res,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "access_key_id", Header: "Access Key ID"},
				{Key: "status", Header: "Status"},
				{Key: "created_at", Header: "Created"},
			},
			Rows: rows,
		},
	}, nil
}
