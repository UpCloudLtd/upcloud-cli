package tokens

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseTokensCommand creates the base 'token' command
func BaseTokensCommand() commands.Command {
	return &tokensCommand{commands.New("token", "Manage tokens")}
}

type tokensCommand struct {
	*commands.BaseCommand
}
