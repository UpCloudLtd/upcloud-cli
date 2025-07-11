package label

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the 'objectstorage label list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List labels for a managed object storage service",
			"upctl object-storage label list <service-uuid>",
			"upctl object-storage label list my-service",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
}

// Execute implements Command.Execute
func (s *listCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	svc := exec.All()

	msg := "Listing labels for object storage service " + serviceUUID
	exec.PushProgressStarted(msg)

	// Get the service to access its labels
	getReq := &request.GetManagedObjectStorageRequest{UUID: serviceUUID}
	service, err := svc.GetManagedObjectStorage(exec.Context(), getReq)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	// Use the existing labels formatting from the labels package
	section := labels.GetLabelsSectionWithResourceType(service.Labels, "managed object storage")
	return output.Combined{section}, nil
}
