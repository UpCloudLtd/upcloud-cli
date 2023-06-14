package servergroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
)

// BaseServergroupCommand creates the base "servergroup" command
func BaseServergroupCommand() commands.Command {
	return &servergroupCommand{
		commands.New("servergroup", "Manage server groups"),
	}
}

type servergroupCommand struct {
	*commands.BaseCommand
}
