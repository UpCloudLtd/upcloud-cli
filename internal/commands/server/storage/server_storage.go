package serverstorage

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
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

func (s *serverStorageCommand) BuildSubCommands(cfg *config.Config) {
	commands.BuildCommand(AttachCommand(), s.Cobra(), cfg)
	commands.BuildCommand(DetachCommand(), s.Cobra(), cfg)
}
