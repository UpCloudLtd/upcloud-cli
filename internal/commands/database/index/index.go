package databaseindex

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseIndexCommand creates the base "database index" command
func BaseIndexCommand() commands.Command {
	return &databaseIndexCommand{
		commands.New("index", "Manage database indices"),
	}
}

type databaseIndexCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (k *databaseIndexCommand) InitCommand() {
	k.Cobra().Aliases = []string{"indices"}
}
