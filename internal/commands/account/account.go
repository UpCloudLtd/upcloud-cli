package account

import (
	"github.com/UpCloudLtd/cli/internal/commands"
)

func AccountCommand() commands.Command {
	return &accountCommand{commands.New("account", "Manage account")}
}

type accountCommand struct {
	*commands.BaseCommand
}
