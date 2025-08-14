package bucket

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the 'objectstorage bucket list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List buckets in a managed object storage service",
			"upctl object-storage bucket list <service-uuid>",
			"upctl object-storage bucket list my-service",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
}

// InitCommand implements Command.InitCommand
func (s *listCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`List buckets in a managed object storage service

Lists all buckets in the specified managed object storage service, showing their names and total size in bytes.`)
}

// Execute implements commands.MultipleArgumentCommand
func (s *listCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Listing buckets in service %s", serviceUUID)
	exec.PushProgressStarted(msg)

	req := &request.GetManagedObjectStorageBucketMetricsRequest{
		ServiceUUID: serviceUUID,
	}

	res, err := svc.GetManagedObjectStorageBucketMetrics(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	rows := []output.TableRow{}
	for _, bucket := range res {
		rows = append(rows, output.TableRow{
			bucket.Name,
			fmt.Sprintf("%d", bucket.TotalSizeBytes),
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: res,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "total_size_bytes", Header: "Total Size Bytes"},
			},
			Rows: rows,
		},
	}, nil
}
