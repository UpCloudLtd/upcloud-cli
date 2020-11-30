package ui

import (
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/cli/internal/validation"
)

func StyleDataTable() table.Style {
	return table.Style{
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
	}
}

func StyleFlagsTable() table.Style {
	return table.Style{
		Name: "FlagsTable",
		Box: table.BoxStyle{
			BottomLeft:       " ",
			BottomRight:      " ",
			BottomSeparator:  " ",
			Left:             " ",
			LeftSeparator:    " ",
			MiddleHorizontal: " ",
			MiddleSeparator:  " ",
			MiddleVertical:   " ",
			PaddingLeft:      "  ",
			PaddingRight:     "",
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
			SeparateRows:    true,
		},
	}
}

func NewDataTable(columnKeys ...string) *DataTable {
	t := &DataTable{columnKeys: columnKeys}
	t.init()
	t.SetStyle(StyleDataTable())
	t.header = table.Row{}
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

func (s *DataTable) SetStyle(style table.Style) {
	s.t.SetStyle(style)
}

func (s *DataTable) SetHeader(hdr table.Row) {
	if hdr != nil && len(hdr) != len(s.columnKeys) {
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
		if len(s.header) == 0 && s.header != nil {
			header = append(header, key)
		} else if s.header != nil {
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
			if _, ok := s.rows[0][pos].(time.Time); ok && cfg.Transformer == nil {
				cfg.Transformer = func(val interface{}) string {
					tv, ok := val.(time.Time)
					if !ok {
						return fmt.Sprintf("%s", val)
					}
					return FormatTime(tv)
				}
			}
			if _, ok := s.rows[0][pos].(float64); ok && cfg.Transformer == nil {
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
		s.t.AppendHeader(header)
	}
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

func TerminalHeight() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.Split(string(out), " ")[0])
}

func (s *DataTable) Paginate(writer io.Writer) error {
	var i int
	var rowCount int
	b := make([]byte, 1)

	fmt.Fprintln(writer)
	for i < len(s.rows) {
		height, err := TerminalHeight()
		if err != nil {
			return err
		}

		if i-len(s.rows) < height && i > 0 {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
		t := NewDataTable(s.columnKeys...)
		t.OverrideColumnKeys(s.overrideColumnKeys...)
		t.SetHeader(s.header)
		for key, value := range s.columnConfig {
			t.SetColumnConfig(key, *value)
		}
		t.SetStyle(*s.t.Style())

		for rowCount < height-4 && i < len(s.rows) {
			t.AppendRow(s.rows[i])
			rowCount++
			i++
		}

		fmt.Fprintln(writer, t.Render())

		if i < len(s.rows) {
			fmt.Println("Press Enter for next page...")
			os.Stdin.Read(b)
			rowCount = 0
			continue
		}
	}
	fmt.Fprintln(writer)
	return nil
}
