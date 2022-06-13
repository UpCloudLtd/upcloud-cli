package database

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
)

// BaseDatabaseCommand creates the base "database" command
func BaseDatabaseCommand() commands.Command {
	return &databaseCommand{
		commands.New("database", "Manage databases"),
	}
}

type databaseCommand struct {
	*commands.BaseCommand
}
