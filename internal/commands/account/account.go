package account

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseAccountCommand creates the base 'account' command
func BaseAccountCommand() commands.Command {
	return &accountCommand{commands.New("account", "Manage accounts")}
}

type accountCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (acc *accountCommand) InitCommand() {
	acc.Cobra().Aliases = []string{"acc"}
}
