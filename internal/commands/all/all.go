package all

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/plan"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func BuildCommands(mainCommand commands.Command, mainConfig *config.Config) {
	c := func() *config.Config { return config.New(mainConfig.Viper()) }
	s := upapi.Service(c())

	// Plans
	planCommand := commands.BuildCommand(plan.PlanCommand(), mainCommand, c())
	commands.BuildCommand(plan.ListCommand(), planCommand, c())

	// Servers
	serverCommand := commands.BuildCommand(server.ServerCommand(), mainCommand, c())
	commands.BuildCommand(server.ListCommand(), serverCommand, c())
	commands.BuildCommand(server.ShowCommand(), serverCommand, c())
	commands.BuildCommand(server.StartCommand(), serverCommand, c())
	commands.BuildCommand(server.StopCommand(), serverCommand, c())
	commands.BuildCommand(server.CreateCommand(), serverCommand, c())
	commands.BuildCommand(server.DeleteCommand(), serverCommand, c())

	// Storages
	storageCommand := commands.BuildCommand(storage.StorageCommand(), mainCommand, c())

	commands.BuildCommand(storage.ListCommand(s), storageCommand, c())
	commands.BuildCommand(storage.ShowCommand(s), storageCommand, c())
	commands.BuildCommand(storage.CreateCommand(s), storageCommand, c())
	commands.BuildCommand(storage.DeleteCommand(s), storageCommand, c())
	commands.BuildCommand(storage.ImportCommand(s), storageCommand, c())
}
