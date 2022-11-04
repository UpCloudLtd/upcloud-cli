package databaseconnection

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
)

// BaseConnectionCommand creates the base "connection" command
func BaseConnectionCommand() commands.Command {
	return &databaseConnectionCommand{
		commands.New("connection", "Manage database connections"),
	}
}

type databaseConnectionCommand struct {
	*commands.BaseCommand
}
