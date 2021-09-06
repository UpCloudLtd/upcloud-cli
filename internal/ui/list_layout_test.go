package ui_test

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
)

func TestShowDetailsRender(t *testing.T) {
	t.Parallel()
	for _, testcase := range []struct {
		name           string
		buildViewFn    func(*ui.ListLayout)
		style          ui.ListLayoutConfig
		expectedOutput string
	}{
		{
			name: "Default style",
			buildViewFn: func(layout *ui.ListLayout) {
				layout.AppendSection("Parent", "Child 1", "Child 2")
			},
			style: ui.ListLayoutDefault,
			expectedOutput: `  
  Parent
    Child 1
    Child 2`,
		},
		{
			name: "Default style with note",
			buildViewFn: func(layout *ui.ListLayout) {
				layout.AppendSectionWithNote("Parent", "Child 1", "(note)")
			},
			style: ui.ListLayoutDefault,
			expectedOutput: `  
  Parent
    Child 1
    
    (note)`,
		},
		{
			name: "Nested table style",
			buildViewFn: func(layout *ui.ListLayout) {
				layout.AppendSection("Parent", "Child 1", "Child 2")
			},
			style: ui.ListLayoutNestedTable,
			expectedOutput: `Parent
  
  Child 1
  
  Child 2`,
		},
		{
			name: "Nested table with note",
			buildViewFn: func(layout *ui.ListLayout) {
				layout.AppendSectionWithNote("Parent", "Child 1", "(note)")
			},
			style: ui.ListLayoutNestedTable,
			expectedOutput: `Parent
  
  Child 1
  
  (note)`,
		},
	} {
		// grab local reference for parallel tests
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			text.DisableColors()
			view := ui.NewListLayout(testcase.style)
			testcase.buildViewFn(view)
			// fmt.Println(view.Render())
			assert.Equal(t, testcase.expectedOutput, view.Render())
		})
	}
}
