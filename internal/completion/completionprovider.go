package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/spf13/cobra"
)

// Provider should be implemented by a command that can provide argument completion
type Provider interface {
	CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective)
}
