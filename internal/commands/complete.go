package commands

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	"github.com/spf13/cobra"
)

type CompleteFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func CompletionFunc(provider completion.Provider, cfg *config.Config) CompleteFunc {
	return func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		svc, err := cfg.CreateService()
		if err != nil {
			return completion.None(toComplete)
		}

		return provider.CompleteArgument(cfg.Context(), svc, toComplete)
	}
}
