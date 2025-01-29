package ui

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/terminal"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func styleDetails(t *table.Table) {
	t.SetStyle(table.Style{
		Name: "Details",
		Box: table.BoxStyle{
			BottomLeft:       "",
			BottomRight:      "",
			BottomSeparator:  "",
			Left:             "",
			LeftSeparator:    "",
			MiddleHorizontal: " ",
			MiddleSeparator:  "",
			MiddleVertical:   "",
			PaddingLeft:      "",
			PaddingRight:     " ",
			Right:            "",
			RightSeparator:   "",
			TopLeft:          "",
			TopRight:         "",
			TopSeparator:     "",
			UnfinishedRow:    "",
		},
		Color: table.ColorOptions{
			Footer:       nil,
			Header:       nil,
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
			SeparateFooter:  true,
			SeparateHeader:  true,
			SeparateRows:    false,
		},
	})
}

// NewDetailsView returns a new output data container for details rendered as a table
func NewDetailsView() *DetailsView {
	t := &DetailsView{t: &table.Table{}}
	styleDetails(t.t)
	return t
}

// DetailsView is an output data container for details rendered as a table
type DetailsView struct {
	t              *table.Table
	rows           []table.Row
	rowTransformer func(row table.Row) table.Row
	headerWidth    int
}

// Render renders a data container for details as a table
func (s *DetailsView) Render() string {
	if len(s.rows) < 1 {
		return ""
	}
	s.t.ResetRows()
	const headerMaxWidth = 20
	if s.headerWidth == 0 {
		s.headerWidth = headerMaxWidth
	}
	widthRemaining := terminal.GetTerminalWidth()
	var colConfigs []table.ColumnConfig
	for i := range s.rows[0] {
		if i < len(s.rows[0])-1 {
			colConfigs = append(colConfigs, table.ColumnConfig{
				Number:    i + 1,
				Align:     text.AlignLeft,
				Colors:    DefaultHeaderColours,
				AutoMerge: true,
				WidthMax:  s.headerWidth,
			})
			widthRemaining -= s.headerWidth
			continue
		}
		colConfigs = append(colConfigs, table.ColumnConfig{WidthMax: widthRemaining})
	}
	s.t.SetColumnConfigs(colConfigs)
	if s.rowTransformer != nil {
		for _, row := range s.rows {
			s.t.AppendRow(s.rowTransformer(row))
		}
	} else {
		s.t.AppendRows(s.rows)
	}
	return s.t.Render()
}

// SetRowSeparators sets whether Render() outputs separators ascii lines between data rows
func (s *DetailsView) SetRowSeparators(v bool) {
	style := s.t.Style()
	style.Options.SeparateRows = v
	style.Box.MiddleHorizontal = "─"
	style.Box.MiddleSeparator = "┼"
	style.Box.MiddleVertical = "│"
}

// SetRowSpacing sets whether Render() outputs separators empty lines between data rows
func (s *DetailsView) SetRowSpacing(v bool) {
	style := s.t.Style()
	style.Options.SeparateRows = v
	style.Box.MiddleHorizontal = " "
	style.Box.MiddleSeparator = "  │"
	style.Box.LeftSeparator = "│"
	style.Box.MiddleVertical = "│"
}

// SetRowTransformer sets a method to transform rwos before rendering
func (s *DetailsView) SetRowTransformer(fn func(row table.Row) table.Row) {
	s.rowTransformer = fn
}

// Append appends new rows to the DetailsView
func (s *DetailsView) Append(rows ...table.Row) {
	s.rows = append(s.rows, rows...)
}

// SetHeaderWidth sets the width of the header rendered
func (s *DetailsView) SetHeaderWidth(width int) {
	s.headerWidth = width
}
