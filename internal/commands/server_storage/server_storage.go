package server_storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
)

const (
	minStorageSize    = 10
	maxServerActions  = 10
	PositionalArgHelp = "<UUID/Title/Hostname...>"
)

func ServerStorageCommand() commands.Command {
	return &serverStorageCommand{commands.New("storage", "Manage server storages")}
}

type serverStorageCommand struct {
	*commands.BaseCommand
}

// A lot of stuff being imported from /commands/server.
// Because storage server was originally refactored out
// from there and makes sense not to duplicate the code here.
