package output

import (
	"encoding/json"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v2"
)

// DetailRow represents a single row in the details view, with a title and a value
type DetailRow struct {
	Title string
	Value interface{}
}

// DetailSection represents a section in the details view
type DetailSection struct {
	Title string
	Rows  []DetailRow
}

// MarshalJSON implements json.Marshaler
func (d DetailSection) MarshalJSON() ([]byte, error) {
	jsonObject := map[string]interface{}{}
	for _, r := range d.Rows {
		jsonObject[r.Title] = r.Value
	}
	return json.Marshal(jsonObject)
}

// Details implements output.Command for a details-style view
type Details struct {
	Sections []DetailSection
}

// MarshalJSON implements json.Marshaler
func (d Details) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Sections)
}

// MarshalYAML marshals details and returns the YAML as []byte
// nb. does *not* implement yaml.Marshaler
func (d Details) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(map[string]interface{}{
		"sections": d.Sections,
	})
}

// MarshalHuman marshals details and returns a human readable []byte
func (d Details) MarshalHuman() ([]byte, error) {
	layout := ui.ListLayoutDefault
	l := ui.NewListLayout(layout)
	for _, sec := range d.Sections {
		dCommon := ui.NewDetailsView()
		// TODO: this logic should prooobably be in the table rendering logic.
		hWidth := 10
		for _, row := range sec.Rows {
			if len(row.Title) > hWidth {
				hWidth = len(row.Title)
			}
		}
		dCommon.SetHeaderWidth(hWidth)
		for _, row := range sec.Rows {
			dCommon.Append(table.Row{row.Title, row.Value})
		}
		l.AppendSection(sec.Title, dCommon.Render())

	}
	return []byte(l.Render()), nil
}
