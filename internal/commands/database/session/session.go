package databasesession

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseSessionCommand creates the base "session" command
func BaseSessionCommand() commands.Command {
	return &databaseSessionCommand{
		commands.New("session", "Manage database sessions"),
	}
}

type databaseSessionCommand struct {
	*commands.BaseCommand
}
