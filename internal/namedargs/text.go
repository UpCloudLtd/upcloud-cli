package namedargs

import (
	"fmt"
	"strings"
)

// Returns description for --zone argument, e.g. "Zone where to create the resource...".
func ZoneDescription(resource string) string {
	return fmt.Sprintf("Zone where to create the %s. Run `upctl zone list` to list all available zones.", resource)
}

// ValidValuesHelp wraps values in backticks and adds human readable separators.
// For example, "`one`, `two` and `three`".
func ValidValuesHelp(values ...string) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		return fmt.Sprintf("`%s`", values[0])
	}

	return fmt.Sprintf(
		"`%s` and `%s`",
		strings.Join(values[:len(values)-1], "`, `"),
		values[len(values)-1],
	)
}
