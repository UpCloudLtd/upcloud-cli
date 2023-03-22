package output

import "github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"

// GetLabelsSection returns labels table as CombinedSection
func GetLabelsSection(labels []upcloud.Label) CombinedSection {
	var rows []TableRow
	for _, i := range labels {
		rows = append(rows, TableRow{i.Key, i.Value})
	}

	return CombinedSection{
		Key:   "labels",
		Title: "Labels:",
		Contents: Table{
			Columns: []TableColumn{
				{Key: "key", Header: "Key"},
				{Key: "value", Header: "Value"},
			},
			Rows: rows,
		},
	}
}

// InsertLabelsIntoCombined returns an Combined output with table generated from given labels list inserted after overview section, if there were labels defined.
func InsertLabelsIntoCombined(combined Combined, labels []upcloud.Label) Combined {
	var rows []TableRow
	for _, i := range labels {
		rows = append(rows, TableRow{i.Key, i.Value})
	}

	r := append(
		combined[:2],
		combined[1:]...,
	)
	r[1] = GetLabelsSection(labels)
	return r
}
