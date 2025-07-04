package label

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// RemoveCommand creates the 'objectstorage label remove' command
func RemoveCommand() commands.Command {
	return &removeCommand{
		BaseCommand: commands.New(
			"remove",
			"Remove labels from a managed object storage service",
			"upctl object-storage label remove <service-uuid> --key env --key team",
			"upctl object-storage label remove my-service --key env",
		),
	}
}

type removeCommand struct {
	*commands.BaseCommand
	labelKeys []string
}

// InitCommand implements Command.InitCommand
func (s *removeCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringArrayVar(&s.labelKeys, "key", nil, "Label keys to remove, multiple can be declared.\nUsage: --key env --key owner")
}

// Execute implements Command.Execute
func (s *removeCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	if len(s.labelKeys) == 0 {
		return nil, fmt.Errorf("at least one label key must be specified")
	}

	svc := exec.All()

	// Get current service to get existing labels
	getReq := &request.GetManagedObjectStorageRequest{UUID: serviceUUID}
	current, err := svc.GetManagedObjectStorage(exec.Context(), getReq)
	if err != nil {
		return commands.HandleError(exec, "Failed to get current service", err)
	}

	// Remove specified labels
	updatedLabels := removeLabels(current.Labels, s.labelKeys)

	// Update the service
	msg := fmt.Sprintf("Removing labels from object storage service %s", serviceUUID)
	exec.PushProgressStarted(msg)

	modifyReq := &request.ModifyManagedObjectStorageRequest{
		UUID:   serviceUUID,
		Labels: &updatedLabels,
	}

	res, err := svc.ModifyManagedObjectStorage(exec.Context(), modifyReq)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
		{Title: "Name", Value: res.Name},
		{Title: "Labels", Value: commands.FormatLabelsCondensed(res.Labels)},
	}}, nil
}

// removeLabels removes labels with specified keys from the existing labels
func removeLabels(existing []upcloud.Label, keysToRemove []string) []upcloud.Label {
	removeSet := make(map[string]bool)
	for _, key := range keysToRemove {
		removeSet[key] = true
	}

	var result []upcloud.Label
	for _, label := range existing {
		if !removeSet[label.Key] {
			result = append(result, label)
		}
	}
	return result
}
