package label

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseLabelCommand creates the base "object-storage label" command
func BaseLabelCommand() commands.Command {
	return &labelCommand{
		BaseCommand: commands.New("label", "Manage labels in managed object storage services"),
	}
}

type labelCommand struct {
	*commands.BaseCommand
}
