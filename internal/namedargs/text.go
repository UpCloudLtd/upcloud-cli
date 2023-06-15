package namedargs

import "fmt"

// Returns description for --zone argument, e.g. "Zone where to create the resource...".
func ZoneDescription(resource string) string {
	return fmt.Sprintf("Zone where to create the %s. Run `upctl zone list` to list all available zones.", resource)
}
