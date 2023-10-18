package output

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
)

// DetailRow represents a single row in the details view, with a title and a value
type DetailRow struct {
	Title  string // used for human-readable representations
	Key    string // user for machine-readable (json, yaml) representations
	Value  interface{}
	Colour text.Colors
	Format func(val interface{}) (text.Colors, string, error)
}

// DetailSection represents a section in the details view
type DetailSection struct {
	Title string // used for human-readable representations
	Key   string // user for machine-readable (json, yaml) representations
	Rows  []DetailRow
}

// MarshalJSON implements json.Marshaler
func (d DetailSection) MarshalJSON() ([]byte, error) {
	jsonObject := map[string]interface{}{}
	for _, r := range d.Rows {
		jsonObject[r.Key] = r.Value
	}
	return json.Marshal(jsonObject)
}

// Details implements output.Output for a details-style view
type Details struct {
	Sections []DetailSection
}

// MarshalJSON implements json.Marshaler
func (d Details) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(mapSections(d.Sections), "", "  ")
}

func mapSections(sections []DetailSection) map[string]interface{} {
	out := make(map[string]interface{})
	for _, section := range sections {
		if section.Key != "" {
			out[section.Key] = mapSectionRows(section.Rows)
		} else {
			for k, v := range mapSectionRows(section.Rows) {
				out[k] = v
			}
		}
	}
	return out
}

func mapSectionRows(rows []DetailRow) map[string]interface{} {
	out := make(map[string]interface{})
	for _, row := range rows {
		out[row.Key] = row.Value
	}
	return out
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
			switch {
			case row.Format != nil:
				colour, formatted, err := row.Format(row.Value)
				if err != nil {
					return nil, fmt.Errorf("error formatting row '%v': %w", row.Key, err)
				}
				dCommon.Append(table.Row{row.Title, colour.Sprintf("%v", formatted)})
			case row.Colour != nil:
				dCommon.Append(table.Row{row.Title, row.Colour.Sprintf("%v", row.Value)})
			default:
				dCommon.Append(table.Row{row.Title, row.Value})
			}
		}
		l.AppendSection(sec.Title, dCommon.Render())
	}

	// add a newline at the end
	return append([]byte(l.Render()), '\n'), nil
}

// MarshalRawMap implements output.Output
func (d Details) MarshalRawMap() (map[string]interface{}, error) {
	return mapSections(d.Sections), nil
}
