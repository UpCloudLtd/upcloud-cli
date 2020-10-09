package table

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"

	"github.com/UpCloudLtd/cli/internal/terminal"
)

var DefaultTableHeaderColours = tablewriter.Colors{tablewriter.Bold}

type Table struct {
	*tablewriter.Table
	header         []string
	headerPos      map[string]int
	visibleColumns []string
	columnKeys     []string
	columnKeyPos   map[string]int
	rows           [][]string
	rowColours     [][]tablewriter.Colors
	curRow         int
}

func (s *Table) AddColumn(value string, colour tablewriter.Colors) {
	if len(s.rows)-1 != s.curRow {
		s.rows = append(s.rows, make([]string, 0, len(s.header)))
		s.rowColours = append(s.rowColours, make([]tablewriter.Colors, 0, len(s.header)))
	}
	s.rows[s.curRow] = append(s.rows[s.curRow], value)
	if len(colour) == 0 || !terminal.Colours() {
		colour = tablewriter.Colors{}
	}
	s.rowColours[s.curRow] = append(s.rowColours[s.curRow], colour)
}

func (s *Table) NextRow() {
	if len(s.rows)-1 == s.curRow && len(s.rows[s.curRow]) > 0 {
		s.curRow += 1
	}
}

func (s *Table) SetVisibleColumns(headersOrColumnKeys ...string) {
	s.visibleColumns = headersOrColumnKeys
}

func (s *Table) SetColumnKeys(keys ...string) {
	s.columnKeys = keys
	for i, k := range keys {
		s.columnKeyPos[k] = i
	}
}

func (s *Table) Render() {
	var visiblePositions []int
	for _, vcol := range s.visibleColumns {
		headerPos, hExists := s.headerPos[vcol]
		columnPos, ckExists := s.columnKeyPos[vcol]
		if !hExists && !ckExists {
			continue
		}
		var pos int
		if hExists {
			pos = headerPos
		}
		if ckExists {
			pos = columnPos
		}
		visiblePositions = append(visiblePositions, pos)
	}
	var header = s.header[:]
	if len(visiblePositions) > 0 {
		header = make([]string, 0, len(visiblePositions))
		for _, pos := range visiblePositions {
			header = append(header, s.header[pos])
		}
	}
	s.SetHeader(header)
	if terminal.Colours() {
		var headerColours []tablewriter.Colors
		for i := 0; i < len(header); i++ {
			headerColours = append(headerColours, DefaultTableHeaderColours)
		}
		s.SetHeaderColor(headerColours...)
	}
	for i, row := range s.rows {
		colours := s.rowColours[i]
		if len(visiblePositions) > 0 {
			newRow := make([]string, 0, len(visiblePositions))
			newColours := make([]tablewriter.Colors, 0, len(visiblePositions))
			for _, pos := range visiblePositions {
				newRow = append(newRow, row[pos])
				newColours = append(newColours, colours[pos])
			}
			row = newRow
			colours = newColours
		}
		s.Table.Rich(row, colours)
	}
	s.Table.Render()
}

func New(output io.Writer, header ...string) (*Table, error) {
	tw := tablewriter.NewWriter(output)
	t := &Table{
		Table:        tw,
		header:       make([]string, 0, len(header)),
		headerPos:    make(map[string]int),
		columnKeyPos: make(map[string]int),
	}
	tw.SetBorder(false)
	tw.SetCenterSeparator(" ")
	tw.SetColumnSeparator(" ")
	tw.SetRowSeparator("─")
	tw.SetAutoFormatHeaders(false)
	for i, hdr := range header {
		t.header = append(t.header, hdr)
		if _, exists := t.headerPos[hdr]; exists {
			return nil, fmt.Errorf("header %q already exists", hdr)
		}
		t.headerPos[hdr] = i
	}
	return t, nil
}

func NewDetails(output io.Writer) *Table {
	t, _ := New(output)
	t.SetColumnSeparator("│")
	t.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	return t
}

func MapColours(header []string, row []string, mapper func(hdr string, row string) tablewriter.Colors) []tablewriter.Colors {
	var r []tablewriter.Colors
	for i := 0; i < len(header); i++ {
		r = append(r, mapper(header[i], row[i]))
	}
	return r
}
