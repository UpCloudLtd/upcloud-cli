package all

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/plan"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/config"
)

func BuildCommands(mainCommand commands.Command, mainConfig *config.Config) {
	// Plans
	planCommand := commands.BuildCommand(plan.PlanCommand(), mainCommand, config.New(mainConfig.Viper()))
	commands.BuildCommand(plan.ListCommand(), planCommand, config.New(mainConfig.Viper()))

	// Servers
	serverCommand := commands.BuildCommand(server.ServerCommand(), mainCommand, config.New(mainConfig.Viper()))
	commands.BuildCommand(server.ListCommand(), serverCommand, config.New(mainConfig.Viper()))
	commands.BuildCommand(server.ShowCommand(), serverCommand, config.New(mainConfig.Viper()))
	commands.BuildCommand(server.StartCommand(), serverCommand, config.New(mainConfig.Viper()))
	commands.BuildCommand(server.StopCommand(), serverCommand, config.New(mainConfig.Viper()))
}
