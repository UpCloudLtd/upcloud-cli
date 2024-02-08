package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// Account implements argument completion for accounts by username.
type Account struct{}

// make sure Account implements the interface
var _ Provider = Account{}

// CompleteArgument implements completion.Provider
func (s Account) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	accounts, err := svc.GetAccountList(ctx)
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range accounts {
		vals = append(vals, v.Username)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
