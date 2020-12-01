package ui

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func StyleDetails(t *table.Table) {
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

func NewDetailsView() *DetailsView {
	t := &DetailsView{t: &table.Table{}}
	StyleDetails(t.t)
	return t
}

type DetailsView struct {
	t              *table.Table
	rows           []table.Row
	rowTransformer func(row table.Row) table.Row
	headerWidth    int
}

func (s *DetailsView) Render() string {
	if len(s.rows) < 1 {
		return ""
	}
	s.t.ResetRows()
	const headerMaxWidth = 20
	if s.headerWidth == 0 {
		s.headerWidth = headerMaxWidth
	}
	widthRemaining := 140
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

func (s *DetailsView) SetRowSeparators(v bool) {
	style := s.t.Style()
	style.Options.SeparateRows = v
	style.Box.MiddleHorizontal = "─"
	style.Box.MiddleSeparator = "┼"
	style.Box.MiddleVertical = "│"
}

func (s *DetailsView) SetRowSpacing(v bool) {
	style := s.t.Style()
	style.Options.SeparateRows = v
	style.Box.MiddleHorizontal = " "
	style.Box.MiddleSeparator = "  │"
	style.Box.LeftSeparator = "│"
	style.Box.MiddleVertical = "│"
}

func (s *DetailsView) SetRowTransformer(fn func(row table.Row) table.Row) {
	s.rowTransformer = fn
}

func (s *DetailsView) AppendRow(row table.Row) {
	s.rows = append(s.rows, row)
}

func (s *DetailsView) AppendRows(rows []table.Row) {
	for _, row := range rows {
		s.AppendRow(row)
	}
}

func (s *DetailsView) SetHeaderWidth(width int) {
	s.headerWidth = width
}
