package output

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/jedib0t/go-pretty/v6/text"
)

// BoolFormat returns val formatted as a boolean
func BoolFormat(val interface{}) (text.Colors, string, error) {
	if boolVal, ok := val.(bool); ok {
		if boolVal {
			return ui.DefaultBooleanColoursTrue, "yes", nil
		}
		return ui.DefaultBooleanColoursFalse, "no", nil
	}
	return nil, "", fmt.Errorf("cannot parse '%v' (%T) as boolean", val, val)
}
