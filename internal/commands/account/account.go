package account

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
)

// BaseAccountCommand creates the base 'account' command
func BaseAccountCommand() commands.Command {
	return &accountCommand{commands.New("account", "Manage account", "")}
}

type accountCommand struct {
	*commands.BaseCommand
}
