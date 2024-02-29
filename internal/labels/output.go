package labels

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// GetLabelsSectionWithResourceType returns labels table as output.CombinedSection with resource type in the empty message.
func GetLabelsSectionWithResourceType(labels []upcloud.Label, resourceType string) output.CombinedSection {
	var rows []output.TableRow
	for _, i := range labels {
		rows = append(rows, output.TableRow{i.Key, i.Value})
	}

	return output.CombinedSection{
		Key:   "labels",
		Title: "Labels:",
		Contents: output.Table{
			Columns: []output.TableColumn{
				{Key: "key", Header: "Key"},
				{Key: "value", Header: "Value"},
			},
			Rows:         rows,
			EmptyMessage: fmt.Sprintf("No labels defined for this %s.", resourceType),
		},
	}
}

// GetLabelsSection returns labels table as output.CombinedSection
func GetLabelsSection(labels []upcloud.Label) output.CombinedSection {
	return GetLabelsSectionWithResourceType(labels, "resource")
}
