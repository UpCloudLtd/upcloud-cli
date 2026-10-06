package labels

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// GetLabelsSectionWithResourceType returns labels table as output.CombinedSection with resource type in the empty message.
func GetLabelsSectionWithResourceType(labels []upcloud.Label, resourceType string) output.CombinedSection {
	return LabelsSection(labels, func(l upcloud.Label) (string, string) { return l.Key, l.Value }, resourceType)
}

// LabelsSection returns labels of any SDK type as output.CombinedSection, using keyValue to read each label.
func LabelsSection[T any](labels []T, keyValue func(T) (key, value string), resourceType string) output.CombinedSection {
	var rows []output.TableRow
	for _, label := range labels {
		key, value := keyValue(label)
		rows = append(rows, output.TableRow{key, value})
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
