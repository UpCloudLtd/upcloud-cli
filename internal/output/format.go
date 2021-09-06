package output

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// BoolFormat returns val formatted as a boolean.
func BoolFormat(val interface{}) (text.Colors, string, error) {
	if vb, ok := val.(bool); ok {
		if vb {
			return ui.DefaultBooleanColoursTrue, "yes", nil
		}

		return ui.DefaultBooleanColoursFalse, "no", nil
	} else if upb, ok := val.(upcloud.Boolean); ok {
		if upb == upcloud.True {
			return ui.DefaultBooleanColoursTrue, "yes", nil
		}

		return ui.DefaultBooleanColoursFalse, "no", nil
	}

	return nil, "", fmt.Errorf("cannot parse '%v' (%T) as boolean", val, val)
}
