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

type ListLayoutConfig struct {
	MarginLeft    bool
	MarginTop     bool
	MarginBottom  bool
	PadTop        bool
	PadBottom     bool
	NoteSeparator bool
}

var ListLayoutDefault = ListLayoutConfig{
	MarginLeft:    true,
	MarginTop:     true,
	MarginBottom:  false,
	PadTop:        false,
	PadBottom:     false,
	NoteSeparator: true,
}

var ListLayoutNestedTable = ListLayoutConfig{
	MarginLeft:    false,
	MarginTop:     false,
	MarginBottom:  false,
	PadTop:        true,
	PadBottom:     false,
	NoteSeparator: true,
}

type ListLayout struct {
	l     list.Writer
	style ListLayoutConfig
}

func NewListLayout(style ListLayoutConfig) *ListLayout {
	l := list.NewWriter()
	l.SetStyle(defaultListStyle)

	return &ListLayout{
		l:     l,
		style: style,
	}
}

func WrapWithListLayout(text string, style ListLayoutConfig) *ListLayout {
	l := NewListLayout(style)
	l.appendSection("", "", []string{text})
	return l
}

func (s *ListLayout) AppendSectionWithNote(title, sectionBody, note string) {
	s.appendSection(title, note, []string{sectionBody})
}
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
		s.l.AppendItem(DefaultNoteColours.Sprintf(note))
	}
	if s.style.MarginBottom {
		s.appendLine()
	}
	if titleAvailable {
		s.l.UnIndent()
	}
}

func (s *ListLayout) Render() string {
	if s.style.MarginLeft {
		return s.l.Render()
	} else {
		return s.removePadLeft()
	}
}

func (s *ListLayout) appendLine() {
	s.l.AppendItem("")
}

func (s *ListLayout) removePadLeft() string {
	// removing the padding from the defaultListStyle caused problems with multi-line items
	// removing the left padding manually with regex
	return regexp.MustCompile("(?m)^ {2}").ReplaceAllString(s.l.Render(), "")
}
