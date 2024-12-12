package partneraccount

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseAccountCommand creates the base "partner account" command
func BaseAccountCommand() commands.Command {
	return &partnerAccountCommand{
		commands.New("account", "Manage accounts associated with partner"),
	}
}

type partnerAccountCommand struct {
	*commands.BaseCommand
}
