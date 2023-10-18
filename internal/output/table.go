package output

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/validation"
)

// TableRow represents a single row of data in a table
type TableRow []interface{}

// TableColumn defines how a particular column is rendered
type TableColumn struct {
	Header string
	Key    string
	Hidden bool
	config *table.ColumnConfig
	Colour text.Colors
	Format func(val interface{}) (text.Colors, string, error)
}

// Table represents command output rendered as a table
type Table struct {
	Columns      []TableColumn
	Rows         []TableRow
	EmptyMessage string
	HideHeader   bool
}

func (s Table) asListOfMaps() []map[string]interface{} {
	jmap := []map[string]interface{}{}
	for _, row := range s.Rows {
		jrow := map[string]interface{}{}
		for i := range row {
			jrow[s.Columns[i].Key] = row[i]
		}
		jmap = append(jmap, jrow)
	}
	return jmap
}

// MarshalJSON implements json.Marshaler
func (s Table) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(s.asListOfMaps(), "", "  ")
}

// MarshalHuman returns table output in a human-readable form
func (s Table) MarshalHuman() ([]byte, error) {
	if len(s.Rows) == 0 && len(s.EmptyMessage) > 0 {
		return []byte(text.FgHiBlack.Sprintf("\n%s\n", s.EmptyMessage)), nil
	}

	t := &table.Table{}
	columnKeyPos := make(map[string]int)
	columnConfig := make(map[string]*table.ColumnConfig)
	for pos, column := range s.Columns {
		columnKeyPos[column.Key] = pos
	}
	t.ResetHeaders()
	t.ResetFooters()
	t.ResetRows()
	t.SetStyle(defaultTableStyle)
	/*
		// TODO: reimplement this if/when necessary
		if len(s.overrideColumnKeys) > 0 {
			columnKeys = s.overrideColumnKeys
		}
	*/
	var header table.Row
	for _, column := range s.Columns {
		pos, ok := columnKeyPos[column.Key]
		if !ok {
			continue
		}
		if column.Header == "" {
			header = append(header, column.Key)
		} else {
			header = append(header, column.Header)
		}
		cfg := column.config
		if cfg == nil {
			cfg = &table.ColumnConfig{}
			column.config = cfg
		}
		cfg.Number = pos + 1
		if len(s.Rows) > 0 {
			// See if the row value can be shows as numeric and if so, align right if not aligned
			if err := validation.Numeric(s.Rows[0][pos]); err == nil && cfg.Align == text.AlignDefault {
				cfg.Align = text.AlignRight
			}
			if _, ok := s.Rows[0][pos].(time.Time); ok && cfg.Transformer == nil {
				cfg.Transformer = func(val interface{}) string {
					tv, ok := val.(time.Time)
					if !ok {
						return fmt.Sprintf("%s", val)
					}
					return ui.FormatTime(tv)
				}
			}
			if _, ok := s.Rows[0][pos].(float64); ok && cfg.Transformer == nil {
				cfg.Transformer = func(val interface{}) string {
					fv, ok := val.(float64)
					if !ok {
						return fmt.Sprintf("%s", val)
					}
					if _, frac := math.Modf(fv); frac != 0 {
						return fmt.Sprintf("%s", val)
					}
					return fmt.Sprintf("%.2f", fv)
				}
			}
		}
	}
	if len(header) > 0 {
		t.AppendHeader(header)
	}
	var columnConfigs []table.ColumnConfig
	for _, cfg := range columnConfig {
		columnConfigs = append(columnConfigs, *cfg)
	}
	t.SetColumnConfigs(columnConfigs)
	for _, row := range s.Rows {
		var arow table.Row
		for _, column := range s.Columns {
			if _, ok := columnKeyPos[column.Key]; !ok {
				continue
			}
			val := row[columnKeyPos[column.Key]]
			if column.Format != nil {
				colour, formatted, err := column.Format(val)
				if err != nil {
					return nil, fmt.Errorf("error formatting column '%v': %w", column.Key, err)
				}
				val = colour.Sprintf("%v", formatted)
			} else if column.Colour != nil {
				val = column.Colour.Sprintf("%v", val)
			}
			arow = append(arow, val)
		}
		if len(arow) > 0 {
			t.AppendRow(arow)
		}
	}
	// print one newline before and after to conform with details view (which is controlled by go-pretty)
	return append([]byte{'\n'}, append([]byte(t.Render()), '\n')...), nil
}

// MarshalRawMap implements output.Output
func (s Table) MarshalRawMap() (map[string]interface{}, error) {
	// TODO: make this better..
	return map[string]interface{}{
		"table": s.asListOfMaps(),
	}, nil
}

var defaultTableStyle = table.Style{
	Name: "DataTable",
	Box: table.BoxStyle{
		BottomLeft:       " ",
		BottomRight:      " ",
		BottomSeparator:  " ",
		Left:             " ",
		LeftSeparator:    " ",
		MiddleHorizontal: "─",
		MiddleSeparator:  " ",
		MiddleVertical:   " ",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		Right:            " ",
		RightSeparator:   " ",
		TopLeft:          " ",
		TopRight:         " ",
		TopSeparator:     " ",
		UnfinishedRow:    " ",
	},
	Color: table.ColorOptions{
		Footer:       ui.DefaultHeaderColours,
		Header:       ui.DefaultHeaderColours,
		Row:          nil,
		RowAlternate: nil,
	},
	Format: table.FormatOptions{
		Footer: text.FormatDefault,
		Header: text.FormatDefault,
		Row:    text.FormatDefault,
	},
	Options: table.Options{
		DrawBorder:      false,
		SeparateColumns: true,
		SeparateFooter:  false,
		SeparateHeader:  true,
		SeparateRows:    false,
	},
}
