package ui

import (
	"fmt"
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestShowDetailsRender(t *testing.T) {
	for _, testcase := range []struct {
		name           string
		buildViewFn    func(*ListLayout)
		style          listLayoutConfig
		expectedOutput string
	}{
		{
			name: "Default style",
			buildViewFn: func(layout *ListLayout) {
				layout.AppendSection("Parent", "Child 1", "Child 2")
			},
			style: ListLayoutDefault,
			expectedOutput: `  
  Parent
    Child 1
    Child 2`,
		},
		{
			name: "Default style with note",
			buildViewFn: func(layout *ListLayout) {
				layout.AppendSectionWithNote("Parent", "Child 1", "(note)")
			},
			style: ListLayoutDefault,
			expectedOutput: `  
  Parent
    Child 1
    
    (note)`,
		},
		{
			name: "Nested table style",
			buildViewFn: func(layout *ListLayout) {
				layout.AppendSection("Parent", "Child 1", "Child 2")
			},
			style: ListLayoutNestedTable,
			expectedOutput: `Parent
  
  Child 1
  
  Child 2`,
		},
		{
			name: "Nested table with note",
			buildViewFn: func(layout *ListLayout) {
				layout.AppendSectionWithNote("Parent", "Child 1", "(note)")
			},
			style: ListLayoutNestedTable,
			expectedOutput: `Parent
  
  Child 1
  
  (note)`,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			text.DisableColors()
			view := NewListLayout(testcase.style)
			testcase.buildViewFn(view)
			fmt.Println(view.Render())
			assert.Equal(t, testcase.expectedOutput, view.Render())
		})
	}
}
