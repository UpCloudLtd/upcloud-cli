package zone

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
)

// BaseZoneCommand creates the base "zone" command
func BaseZoneCommand() commands.Command {
	return &zoneCommand{
		commands.New("zone", "Display zone information"),
	}
}

type zoneCommand struct {
	*commands.BaseCommand
}
