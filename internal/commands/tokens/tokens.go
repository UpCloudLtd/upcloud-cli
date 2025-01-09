package tokens

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseTokensCommand creates the base 'tokens' command
func BaseTokensCommand() commands.Command {
	return &tokensCommand{commands.New("tokens", "Manage tokens")}
}

type tokensCommand struct {
	*commands.BaseCommand
}
