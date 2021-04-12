package serverstorage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
)

const (
	maxServerStorageActions = 10
	// TODO: re-add
	// positionalArgHelp       = "<UUID/Title/Hostname...>"
)

// BaseServerStorageCommand creates the base "server storage" command
func BaseServerStorageCommand() commands.Command {
	return &serverStorageCommand{commands.New("storage", "Manage server storages")}
}

type serverStorageCommand struct {
	*commands.BaseCommand
}
