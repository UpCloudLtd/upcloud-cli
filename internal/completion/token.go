package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// Token implements argument completion for tokens, by id.
type Token struct{}

// make sure Token implements the interface
var _ Provider = Token{}

// CompleteArgument implements completion.Provider
func (s Token) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	tokens, err := svc.GetTokens(ctx, &request.GetTokensRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, t := range *tokens {
		vals = append(vals, t.ID, t.Name)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
