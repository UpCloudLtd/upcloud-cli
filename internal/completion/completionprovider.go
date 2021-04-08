package completion

import (
	"github.com/UpCloudLtd/cli/internal/service"
	"github.com/spf13/cobra"
)

// Completer is the simplest form of completion function
type Completer func(toComplete string) ([]string, cobra.ShellCompDirective)

// Provider should be implemented by a command that can provide argument completion
type Provider interface {
	Generate(services service.AllServices) (Completer, error)
}
