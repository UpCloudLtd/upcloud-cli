package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/spf13/cobra"
)

// Username implements argument completion for zones by id.
type Username struct{}

// make sure Kubernetes implements the interface
var _ Provider = Username{}

// CompleteArgument implements completion.Provider
func (s Username) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	accounts, err := svc.GetAccountList(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, account := range accounts {
		vals = append(vals, account.Username)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
