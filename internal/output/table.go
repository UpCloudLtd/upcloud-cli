package output

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/validation"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"gopkg.in/yaml.v2"
	"math"
	"time"
)

// TableRow represents a single row of data in a table
type TableRow []interface{}

// Table represents command output rendered as a table
type Table struct {
	Headers []string
	Keys    []string
	Visible []string
	Rows    []TableRow
}

func (s Table) asListOfMaps() []map[string]interface{} {
	jmap := []map[string]interface{}{}
	for _, row := range s.Rows {
		jrow := map[string]interface{}{}
		for i := range row {
			jrow[s.Keys[i]] = row[i]
		}
		jmap = append(jmap, jrow)
	}
	return jmap
}

// MarshalJSON implements json.Marshaler
func (s Table) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.asListOfMaps())
}

// MarshalYAML returns table output marshaled to YAML.
func (s Table) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s.asListOfMaps())
}

// MarshalHuman returns table output in a human-readable form
func (s Table) MarshalHuman() ([]byte, error) {
	t := &table.Table{}
	columnKeyPos := make(map[string]int)
	columnConfig := make(map[string]*table.ColumnConfig)
	for pos, key := range s.Keys {
		columnKeyPos[key] = pos
	}
	t.ResetHeaders()
	t.ResetFooters()
	t.ResetRows()
	columnKeys := s.Keys
	/*
		// TODO: reimplement this if/when necessary
		if len(s.overrideColumnKeys) > 0 {
			columnKeys = s.overrideColumnKeys
		}
	*/
	var header table.Row
	for _, key := range columnKeys {
		pos, ok := columnKeyPos[key]
		if !ok {
			continue
		}
		if len(header) == 0 && s.Headers != nil {
			header = append(header, key)
		} else if s.Headers != nil {
			header = append(header, s.Headers[pos])
		}
		cfg, ok := columnConfig[key]
		if !ok {
			cfg = &table.ColumnConfig{}
			columnConfig[key] = cfg
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
		for _, key := range columnKeys {
			if _, ok := columnKeyPos[key]; !ok {
				continue
			}
			arow = append(arow, row[columnKeyPos[key]])
		}
		if len(arow) > 0 {
			t.AppendRow(arow)
		}
	}
	return []byte(t.Render()), nil
}
