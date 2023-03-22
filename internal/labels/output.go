package labels

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
)

// GetLabelsSection returns labels table as output.CombinedSection
func GetLabelsSection(labels []upcloud.Label) output.CombinedSection {
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
			EmptyMessage: "No labels defined for this resource.",
		},
	}
}
