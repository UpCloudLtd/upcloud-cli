package completion

import (
	"context"
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/spf13/cobra"
)

// Zone implements argument completion for zones by id.
type Zone struct{}

// make sure Kubernetes implements the interface
var _ Provider = Zone{}

// CompleteArgument implements completion.Provider
func (s Zone) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	zones, err := svc.GetZones(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, zone := range zones.Zones {
		vals = append(vals, fmt.Sprintf("%s\t%s", zone.ID, zone.Description))
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
