package serverstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

const (
	maxServerStorageActions = 10
)

// BaseServerStorageCommand creates the base "server storage" command
func BaseServerStorageCommand() commands.Command {
	return &serverStorageCommand{commands.New("storage", "Manage server storages")}
}

type serverStorageCommand struct {
	*commands.BaseCommand
}
