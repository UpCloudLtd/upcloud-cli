package namedargs

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/spf13/cobra"
)

type CompleteFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

// CompletionFunc creates a flag completion function from given completion provider and config to be passed to Cobra via Command.RegisterFlagCompletionFunc
func CompletionFunc(provider completion.Provider, cfg *config.Config) CompleteFunc {
	return func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		svc, err := cfg.CreateService()
		if err != nil {
			return completion.None(toComplete)
		}

		return provider.CompleteArgument(cfg.Context(), svc, toComplete)
	}
}
