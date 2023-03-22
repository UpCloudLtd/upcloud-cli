package account

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
)

// BaseAccountCommand creates the base 'account' command
func BaseAccountCommand() commands.Command {
	return &accountCommand{commands.New("account", "Manage accounts")}
}

type accountCommand struct {
	*commands.BaseCommand
}
