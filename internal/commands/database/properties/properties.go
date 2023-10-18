package databaseproperties

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// PropertiesCommand creates the "database properties" command
func PropertiesCommand() commands.Command {
	return &propertiesCommand{
		BaseCommand: commands.New("properties", "List available properties for given database type", "upctl database properties pg", "upctl database properties mysql"),
	}
}

type propertiesCommand struct {
	*commands.BaseCommand
}
