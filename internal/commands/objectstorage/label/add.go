package label

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// AddCommand creates the 'objectstorage label add' command
func AddCommand() commands.Command {
	return &addCommand{
		BaseCommand: commands.New(
			"add",
			"Add labels to a managed object storage service",
			"upctl object-storage label add <service-uuid> --label env=production --label team=backend",
			"upctl object-storage label add my-service --label env=production",
		),
	}
}

type addCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	labelStrings []string
}

// InitCommand implements Command.InitCommand
func (s *addCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringArrayVar(&s.labelStrings, "label", nil, "Labels to add in `key=value` format, multiple can be declared.\nUsage: --label env=dev --label owner=operations")
}

// Execute implements Command.Execute
func (s *addCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	if len(s.labelStrings) == 0 {
		return nil, fmt.Errorf("at least one label must be specified")
	}

	// Parse labels
	newLabels, err := labels.StringsToSliceOfLabels(s.labelStrings)
	if err != nil {
		return nil, err
	}

	svc := exec.All()

	// Get current service to merge labels
	getReq := &request.GetManagedObjectStorageRequest{UUID: serviceUUID}
	current, err := svc.GetManagedObjectStorage(exec.Context(), getReq)
	if err != nil {
		return commands.HandleError(exec, "Failed to get current service", err)
	}

	// Merge existing labels with new ones
	mergedLabels := mergeLabels(current.Labels, newLabels)

	// Update the service
	msg := fmt.Sprintf("Adding labels to object storage service %s", serviceUUID)
	exec.PushProgressStarted(msg)

	modifyReq := &request.ModifyManagedObjectStorageRequest{
		UUID:   serviceUUID,
		Labels: &mergedLabels,
	}
	// TODO: update to use dedicated labels endpoint when available
	res, err := svc.ModifyManagedObjectStorage(exec.Context(), modifyReq)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	combined := output.Combined{
		output.CombinedSection{
			Contents: output.Details{
				Sections: []output.DetailSection{
					{
						Title: "Overview:",
						Rows: []output.DetailRow{
							{Title: "UUID:", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
							{Title: "Name:", Value: res.Name},
						},
					},
				},
			},
		},
		labels.GetLabelsSectionWithResourceType(res.Labels, "managed object storage"),
	}

	return output.MarshaledWithHumanOutput{
		Value:  res,
		Output: combined,
	}, nil
}

// mergeLabels merges existing labels with new labels, with new labels taking precedence
func mergeLabels(existing []upcloud.Label, newLabels []upcloud.Label) []upcloud.Label {
	labelMap := make(map[string]string)

	// Add existing labels
	for _, label := range existing {
		labelMap[label.Key] = label.Value
	}

	// Add/override with new labels
	for _, label := range newLabels {
		labelMap[label.Key] = label.Value
	}

	// Convert back to slice
	var result []upcloud.Label
	for key, value := range labelMap {
		result = append(result, upcloud.Label{Key: key, Value: value})
	}
	return result
}
