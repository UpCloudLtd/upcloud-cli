package token

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseTokenCommand creates the base 'token' command
func BaseTokenCommand() commands.Command {
	return &tokensCommand{commands.New("token", "Manage tokens (EXPERIMENTAL)")}
}

type tokensCommand struct {
	*commands.BaseCommand
}
