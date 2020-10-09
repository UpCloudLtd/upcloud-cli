package ui

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/cli/internal/validation"
)

func StyleDataTable(t *table.Table) {
	t.SetStyle(table.Style{
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
			PaddingLeft:      "  ",
			PaddingRight:     "  ",
			Right:            " ",
			RightSeparator:   " ",
			TopLeft:          " ",
			TopRight:         " ",
			TopSeparator:     " ",
			UnfinishedRow:    " ",
		},
		Color: table.ColorOptions{
			Footer:       DefaultHeaderColours,
			Header:       DefaultHeaderColours,
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
	})
}

func NewDataTable(columnKeys ...string) *DataTable {
	t := &DataTable{columnKeys: columnKeys}
	t.init()
	StyleDataTable(t.t)
	return t
}

type DataTable struct {
	t                  *table.Table
	header             table.Row
	overrideColumnKeys []string
	columnKeys         []string
	columnKeyPos       map[string]int
	columnConfig       map[string]*table.ColumnConfig
	rows               []table.Row
}

func (s *DataTable) init() {
	s.t = &table.Table{}
	s.columnKeyPos = make(map[string]int)
	s.columnConfig = make(map[string]*table.ColumnConfig)
	for pos, key := range s.columnKeys {
		s.columnKeyPos[key] = pos
	}
}

func (s *DataTable) SetHeader(hdr table.Row) {
	if len(hdr) != len(s.columnKeys) {
		panic("uneven number of columns and headers")
	}
	s.header = hdr
}

// Overrides column visibility and order
func (s *DataTable) OverrideColumnKeys(keys ...string) {
	if len(keys) == 0 {
		return
	}
	s.overrideColumnKeys = keys
}

func (s *DataTable) SetColumnConfig(key string, config table.ColumnConfig) {
	if _, ok := s.columnKeyPos[key]; !ok {
		panic(fmt.Sprintf("undeclared column key %q", key))
	}
	s.columnConfig[key] = &config
}

func (s *DataTable) Render() string {
	s.t.ResetHeaders()
	s.t.ResetFooters()
	s.t.ResetRows()
	columnKeys := s.columnKeys
	if len(s.overrideColumnKeys) > 0 {
		columnKeys = s.overrideColumnKeys
	}
	var header table.Row
	for _, cfg := range s.columnConfig {
		cfg.Number = 0
	}
	for _, key := range columnKeys {
		pos, ok := s.columnKeyPos[key]
		if !ok {
			continue
		}
		if len(s.header) == 0 {
			header = append(header, key)
		} else {
			header = append(header, s.header[pos])
		}
		cfg, ok := s.columnConfig[key]
		if !ok {
			cfg = &table.ColumnConfig{}
			s.columnConfig[key] = cfg
		}
		cfg.Number = pos + 1
		if len(s.rows) > 0 {
			// See if the row value can be shows as numeric and if so, align right if not aligned
			if err := validation.Numeric(s.rows[0][pos]); err == nil && cfg.Align == text.AlignDefault {
				cfg.Align = text.AlignRight
			}
		}
	}
	s.t.AppendHeader(header)
	var columnConfigs []table.ColumnConfig
	for _, cfg := range s.columnConfig {
		columnConfigs = append(columnConfigs, *cfg)
	}
	s.t.SetColumnConfigs(columnConfigs)
	for _, row := range s.rows {
		var arow table.Row
		for _, key := range columnKeys {
			if _, ok := s.columnKeyPos[key]; !ok {
				continue
			}
			arow = append(arow, row[s.columnKeyPos[key]])
		}
		if len(arow) > 0 {
			s.t.AppendRow(arow)
		}
	}
	return s.t.Render()
}

func (s *DataTable) AppendRow(row table.Row) {
	if len(row) != len(s.columnKeys) {
		panic("uneven number of columns in a row vs the number of column keys")
	}
	s.rows = append(s.rows, row)
}

func (s *DataTable) AppendRows(rows []table.Row) {
	for _, row := range rows {
		s.AppendRow(row)
	}
}
