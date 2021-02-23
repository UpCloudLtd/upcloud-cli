package serverstorage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
)

const (
	maxServerActions  = 10
	positionalArgHelp = "<UUID/Title/Hostname...>"
)

// BaseServerStorageCommand creates the base "server storage" command
func BaseServerStorageCommand() commands.Command {
	return &serverStorageCommand{commands.New("storage", "Manage server storages")}
}

type serverStorageCommand struct {
	*commands.BaseCommand
}

// A lot of stuff being imported from /commands/server.
// Because storage server was originally refactored out
// from there and makes sense not to duplicate the code here.
// TODO: is this still relevant???
