package billing

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseBillingCommand creates the base 'billing' command
func BaseBillingCommand() commands.Command {
	return &billingCommand{commands.New("billing", "Manage billing and view cost summaries")}
}

type billingCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (b *billingCommand) InitCommand() {
	b.Cobra().Aliases = []string{"bill"}
}