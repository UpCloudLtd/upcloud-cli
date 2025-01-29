package ui

import (
	"regexp"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/text"
)

var defaultListStyle = list.Style{
	Format:           text.FormatDefault,
	CharItemSingle:   " ",
	CharItemTop:      " ",
	CharItemFirst:    " ",
	CharItemMiddle:   " ",
	CharItemVertical: "  ",
	CharItemBottom:   " ",
	CharNewline:      "\n",
	LinePrefix:       "",
	Name:             "ListDefault",
}

type listLayoutConfig struct {
	MarginLeft    bool
	MarginTop     bool
	MarginBottom  bool
	PadTop        bool
	PadBottom     bool
	NoteSeparator bool
}

// ListLayoutDefault defines the default config used for rendering lists
var ListLayoutDefault = listLayoutConfig{
	MarginLeft:    true,
	MarginTop:     true,
	MarginBottom:  false,
	PadTop:        false,
	PadBottom:     false,
	NoteSeparator: true,
}

// ListLayoutNestedTable defines the configuration used for rendering nested tables
var ListLayoutNestedTable = listLayoutConfig{
	MarginLeft:    false,
	MarginTop:     false,
	MarginBottom:  false,
	PadTop:        true,
	PadBottom:     false,
	NoteSeparator: true,
}

// ListLayout is a renderer of list data
type ListLayout struct {
	l     list.Writer
	style listLayoutConfig
}

// NewListLayout returns a new list data renderer
func NewListLayout(style listLayoutConfig) *ListLayout {
	l := list.NewWriter()
	l.SetStyle(defaultListStyle)

	return &ListLayout{
		l:     l,
		style: style,
	}
}

// WrapWithListLayout returns a list data renderer wrapping given text
func WrapWithListLayout(text string, style listLayoutConfig) *ListLayout {
	l := NewListLayout(style)
	l.appendSection("", "", []string{text})
	return l
}

// AppendSectionWithNote appends a section with a note to a list
func (s *ListLayout) AppendSectionWithNote(title, sectionBody, note string) {
	s.appendSection(title, note, []string{sectionBody})
}

// AppendSection appends a section to a list
func (s *ListLayout) AppendSection(title string, sectionBody ...string) {
	s.appendSection(title, "", sectionBody)
}

func (s *ListLayout) appendSection(title, note string, sectionBody []string) {
	if s.style.MarginTop {
		s.appendLine()
	}
	titleAvailable := title != ""

	if titleAvailable {
		s.l.AppendItem(DefaultHeaderColours.Sprint(title))
		s.l.Indent()
	}
	for item := range sectionBody {
		if s.style.PadTop {
			s.appendLine()
		}
		s.l.AppendItem(sectionBody[item])
		if s.style.PadBottom {
			s.appendLine()
		}
	}
	if note != "" {
		if s.style.NoteSeparator {
			s.appendLine()
		}
		s.l.AppendItem(DefaultNoteColours.Sprint(note))
	}
	if s.style.MarginBottom {
		s.appendLine()
	}
	if titleAvailable {
		s.l.UnIndent()
	}
}

// Render renders the ListLayout as configured
func (s *ListLayout) Render() string {
	if s.style.MarginLeft {
		return s.l.Render()
	}
	return s.renderWithoutLeftPadding()
}

func (s *ListLayout) appendLine() {
	s.l.AppendItem("")
}

func (s *ListLayout) renderWithoutLeftPadding() string {
	// removing the padding from the defaultListStyle caused problems with multi-line items
	// removing the left padding manually with regex
	return regexp.MustCompile("(?m)^ {2}").ReplaceAllString(s.l.Render(), "")
}
