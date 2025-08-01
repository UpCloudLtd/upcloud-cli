package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// TimeZone implements argument completion for time zones.
type TimeZone struct{}

// make sure Token implements the interface
var _ Provider = TimeZone{}

// CompleteArgument implements completion.Provider
func (s TimeZone) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	tzs, err := svc.GetTimeZones(ctx)
	if err != nil {
		return None(toComplete)
	}
	return MatchStringPrefix(tzs.TimeZones, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
