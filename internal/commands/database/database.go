package database

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
)

// BaseZoneCommand creates the base "zone" command
func BaseDatabaseCommand() commands.Command {
	return &zoneCommand{
		commands.New("database", "Manage databases"),
	}
}

type zoneCommand struct {
	*commands.BaseCommand
}
