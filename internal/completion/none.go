package completion

import "github.com/spf13/cobra"

// None is a fallback with no completion, for error cases etc.
func None(_ string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}
