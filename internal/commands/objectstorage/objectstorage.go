package objectstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseobjectstorageCommand creates the base "object-storage" command
func BaseobjectstorageCommand() commands.Command {
	baseCmd := commands.New("object-storage", "Manage managed object storage services")
	baseCmd.SetDeprecatedAliases([]string{"objectstorage", "objsto"})

	return &objectstorageCommand{
		BaseCommand: baseCmd,
	}
}

type objectstorageCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (c *objectstorageCommand) InitCommand() {
	c.Cobra().Aliases = []string{"obs", "objectstorage", "objsto"}

	// Deprecating objectstorage and objsto in favour of object-storage
	// TODO: Remove this in the future
	commands.SetDeprecationHelp(c.Cobra(), c.DeprecatedAliases())
}
